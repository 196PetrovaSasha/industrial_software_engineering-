package chapter

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestGetChaptersByUserIdHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           []byte
		expectedStatus int
		mockBehavior   func(mock sql.DB)
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Empty body",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid JSON",
			method:         http.MethodPost,
			body:           []byte(`{"user_id": "invalid json`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid user ID",
			method:         http.MethodPost,
			body:           []byte(`{"user_id": "not a number"}`),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := GetChaptersByUserIdHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewReader(tt.body)

			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			req := httptest.NewRequest(tt.method, "/chapters", reqBody)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetChaptersByUserIdHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := GetChaptersByUserIdHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{"user_id": "123"}`)

	// Настраиваем поведение мока для успешного получения глав
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chapters" WHERE author_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "start_node", "nodes", "characters", "status", "author_id",
		}).AddRow(1, "Chapter 1", 1, "[1,2,3]", "[1,2]", 1, 123))

	req := httptest.NewRequest(http.MethodPost, "/chapters", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	chapters, ok := response["chapters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, chapters, 1)

	chapter := chapters[0].(map[string]interface{})
	assert.Equal(t, "1", chapter["id"])
	assert.Equal(t, "Chapter 1", chapter["name"])
	assert.Equal(t, "1", chapter["start_node"])
	assert.Equal(t, []interface{}{"1", "2", "3"}, chapter["nodes"])
	assert.Equal(t, []interface{}{"1", "2"}, chapter["characters"])
	assert.Equal(t, float64(1), chapter["status"])
	assert.Equal(t, "123", chapter["author"])
}

func TestGetChaptersByUserIdHandler_GetChaptersError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := GetChaptersByUserIdHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{"user_id": "123"}`)

	// Настраиваем поведение мока для ошибки получения глав
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "chapters" WHERE author_id = $1`)).
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest(http.MethodPost, "/chapters", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
