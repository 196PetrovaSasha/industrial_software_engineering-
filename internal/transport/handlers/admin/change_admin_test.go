package admin

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeAdminHandler(t *testing.T) {
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
			body:           []byte("{invalid json"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := AdminAuthorisationHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewReader(tt.body)

			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			req := httptest.NewRequest(tt.method, "/admin/change", reqBody)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestChangeAdminHandler_ChangeAdminSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := AdminAuthorisationHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{
        "id": 1,
        "name": "Updated User",
        "email": "updated@example.com",
        "password": "new_password",
        "admin_status": 1,
        "created_chapters": [1, 2, 3]
    }`)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "admin_status"}).
			AddRow(1, "Old User", "old@example.com", 0))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET name = $1, email = $2, password = $3, admin_status = $4, created_chapters = $5 WHERE id = $6`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/admin/change", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), response["id"])
}

func TestChangeAdminHandler_ChangeAdminError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := AdminAuthorisationHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
	)

	reqBody := []byte(`{"id": 1, "name": "Updated User"}`)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	req := httptest.NewRequest(http.MethodPost, "/admin/change", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()

	handler(w, req)

	assert.Equal(t, 400, w.Code)
}
