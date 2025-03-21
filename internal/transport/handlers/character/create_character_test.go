package character

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

func TestCreateCharacterHandler(t *testing.T) {
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
			body:           []byte(`{"name": "invalid json`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid character creation",
			method:         http.MethodPost,
			body:           []byte(`{"name": "Test Character", "slug": "test-character"`),
			expectedStatus: 400,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	logger := zerolog.Nop()

	handler := CreateCharacterHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		&logger,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewReader(tt.body)

			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			req := httptest.NewRequest(tt.method, "/characters", reqBody)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
