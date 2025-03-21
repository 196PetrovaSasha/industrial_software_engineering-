package migrator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitConfiguration(t *testing.T) {
	// Создаем временный файл .env для теста
	err := os.WriteFile(".env", []byte(`
DB_HOST=localhost
DB_USER=testuser
DB_PASSWORD=testpass
DB_NAME=testdb
DB_PORT=5432
DB_SSLMODE=disable
`), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(".env")

	// Загружаем конфигурацию
	config := InitConfiguration()

	// Проверяем все значения
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"DB_HOST", DB_HOST, "localhost"},
		{"DB_USER", DB_USER, "testuser"},
		{"DB_PASSWORD", DB_PASSWORD, "testpass"},
		{"DB_NAME", DB_NAME, "testdb"},
		{"DB_PORT", DB_PORT, "5432"},
		{"DB_SSLMODE", DB_SSLMODE, "disable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, config[tt.key])
		})
	}
}
