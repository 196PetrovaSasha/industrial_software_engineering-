package chapter

import (
	"bytes"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestUpdateChapterHandler(t *testing.T) {
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
			body:           []byte(`{"id": "invalid json"`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid chapter ID",
			method:         http.MethodPost,
			body:           []byte(`{"id": "not a number"}`),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid node ID",
			method:         http.MethodPost,
			body:           []byte(`{"id": "1", "nodes": ["not a number"]}`),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid character ID",
			method:         http.MethodPost,
			body:           []byte(`{"id": "1", "characters": ["not a number"]}`),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid author ID",
			method:         http.MethodPost,
			body:           []byte(`{"id": "1", "update_author_id": "not a number"}`),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid start node",
			method:         http.MethodPost,
			body:           []byte(`{"id": "1", "start_node": "not a number"}`),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid status",
			method:         http.MethodPost,
			body:           []byte(`{"id": "1", "status": "not a number"}`),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := UpdateChapterHandler(
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

func TestUpdateChapterHandler_UpdateChapterSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := UpdateChapterHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{
        "id": "1",
        "name": "Updated Chapter",
        "nodes": ["1", "2", "3"],
        "characters": ["1", "2"],
        "status": "1",
        "update_author_id": "123",
        "start_node": "1"
    }`)

	// Настраиваем поведение мока для успешного обновления главы
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapters"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapter_nodes"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapter_characters"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodPost, "/chapters", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, 400, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateChapterHandler_UpdateChapterError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := UpdateChapterHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{
        "id": "1",
        "name": "Updated Chapter",
        "nodes": ["1", "2", "3"],
        "characters": ["1", "2"],
        "status": "1",
        "update_author_id": "123",
        "start_node": "1"
    }`)

	// Настраиваем поведение мока для ошибки обновления главы
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapters"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapter_nodes"`)).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	req := httptest.NewRequest(http.MethodPost, "/chapters", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, 400, w.Code)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
