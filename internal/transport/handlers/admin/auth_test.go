package admin

import (
	"bytes"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminAuthorisationHandler(t *testing.T) {
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
		{
			name:           "Empty fields",
			method:         http.MethodPost,
			body:           []byte(`{"email":"","password":""}`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Создаем mock базы данных
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	handler := AdminAuthorisationHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		new(zerolog.Logger),
		AuthConfig{},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewReader(tt.body)

			// Если указано поведение мока, настраиваем его
			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			req := httptest.NewRequest(tt.method, "/auth", reqBody)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
