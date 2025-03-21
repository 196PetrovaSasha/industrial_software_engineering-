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
