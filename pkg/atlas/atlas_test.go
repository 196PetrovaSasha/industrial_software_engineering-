package atlas

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartAtlasSchemaValidation(t *testing.T) {
	tests := []struct {
		name           string
		mockModels     []interface{}
		expectedSQL    string
		expectedReturn bool
	}{
		{
			name: "valid schema generation",
			mockModels: []interface{}{
				&MockMedia{},
				&MockRequest{},
				&MockCharacter{},
				&MockNode{},
				&MockAdmin{},
				&MockChapter{},
				&MockPlayer{},
			},
			expectedSQL:    "CREATE TABLE \"media\" (\"id\" bigserial",
			expectedReturn: true,
		},
		{
			name:           "empty models list",
			mockModels:     []interface{}{},
			expectedSQL:    "",
			expectedReturn: true,
		},
		{
			name: "invalid model type",
			mockModels: []interface{}{
				"invalid_model",
			},
			expectedReturn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем временный файл для вывода
			tmpFile, err := os.CreateTemp("", "atlas-schema-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			// Сохраняем текущий stdout
			oldStdout := os.Stdout
			os.Stdout = tmpFile

			// Вызываем тестируемую функцию
			result := StartAtlasSchemaValidation()

			// Восстанавливаем stdout
			os.Stdout = oldStdout

			// Читаем содержимое файла
			content, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			// Проверяем результат
			assert.Equal(t, tt.expectedReturn, result)

			// Проверяем вывод только для валидного случая
			if tt.expectedSQL != "" {
				assert.Contains(t, string(content), tt.expectedSQL)
			}
		})
	}
}

// MockMedia реализует интерфейс для тестирования
type MockMedia struct{}

func (m *MockMedia) TableName() string {
	return "media"
}

// MockRequest реализует интерфейс для тестирования
type MockRequest struct{}

func (m *MockRequest) TableName() string {
	return "requests"
}

// MockCharacter реализует интерфейс для тестирования
type MockCharacter struct{}

func (m *MockCharacter) TableName() string {
	return "characters"
}

// MockNode реализует интерфейс для тестирования
type MockNode struct{}

func (m *MockNode) TableName() string {
	return "nodes"
}

// MockAdmin реализует интерфейс для тестирования
type MockAdmin struct{}

func (m *MockAdmin) TableName() string {
	return "admins"
}

// MockChapter реализует интерфейс для тестирования
type MockChapter struct{}

func (m *MockChapter) TableName() string {
	return "chapters"
}

// MockPlayer реализует интерфейс для тестирования
type MockPlayer struct{}

func (m *MockPlayer) TableName() string {
	return "players"
}
