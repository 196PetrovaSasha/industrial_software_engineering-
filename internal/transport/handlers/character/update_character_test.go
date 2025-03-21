package character

import (
	"bytes"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
)

func TestUpdateCharacterHandler(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		requestBody      string
		expectedStatus   int
		mockBehavior     func(mock sql.DB)
		expectedBody     map[string]interface{}
		expectedLogCalls []string
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
			name:           "Invalid JSON body",
			method:         http.MethodPost,
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid character ID",
			method:         http.MethodPost,
			requestBody:    `{"id":"invalid","name":"Test Character","slug":"test-character","color":"#FF0000","emotions":{"1":"2"}}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid emotion index",
			method:         http.MethodPost,
			requestBody:    `{"id":"1","name":"Test Character","slug":"test-character","color":"#FF0000","emotions":{"invalid":"2"}}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовый логгер с буфером для проверки сообщений
	var logBuffer bytes.Buffer
	logger := zerolog.New(&logBuffer)

	handler := UpdateCharacterHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		&logger,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/characters", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if len(tt.expectedLogCalls) > 0 {
				lines := strings.Split(logBuffer.String(), "\n")
				assert.Len(t, lines, len(tt.expectedLogCalls)+1) // +1 для пустой последней строки
				for i, expectedCall := range tt.expectedLogCalls {
					assert.Contains(t, lines[i], expectedCall)
				}
				logBuffer.Reset() // Очищаем буфер для следующего теста
			}
		})
	}
}
