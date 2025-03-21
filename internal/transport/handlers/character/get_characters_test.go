package character

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"vn/internal/models"
)

func TestGetCharacterHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		mockBehavior   func(mock sql.DB)
		expectedBody   map[string]interface{}
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	logger := zerolog.Nop()

	handler := GetCharacterHandler(
		&gorm.DB{Config: &gorm.Config{ConnPool: db}},
		&logger,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/characters", nil)
			w := httptest.NewRecorder()

			if tt.mockBehavior != nil {
				tt.mockBehavior(*db)
			}

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestPrepareCharacterForResponse(t *testing.T) {
	tests := []struct {
		name  string
		input *[]models.Character
		want  []ResponseCharacter
	}{
		{
			name:  "Empty slice",
			input: &[]models.Character{},
			want:  nil,
		},
		{
			name: "Single character",
			input: &[]models.Character{
				{
					Id:       1,
					Name:     "test_name",
					Slug:     "test_slug",
					Color:    "#FF0000",
					Emotions: map[int64]int64{1: 2},
				},
			},
			want: []ResponseCharacter{
				{
					Id:       "1",
					Name:     "test_name",
					Slug:     "test_slug",
					Color:    "#FF0000",
					Emotions: map[string]string{"1": "2"},
				},
			},
		},
		{
			name:  "Nil input",
			input: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrepareCharacterForResponse(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
