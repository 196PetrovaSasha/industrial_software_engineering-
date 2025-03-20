package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
)

func TestInitConfiguration(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(DB_HOST, "localhost")
	t.Setenv(DB_USER, "test_user")
	t.Setenv(DB_PASSWORD, "test_password")
	t.Setenv(DB_NAME, "test_db")
	t.Setenv(DB_PORT, "5432")
	t.Setenv(DB_SSLMODE, "disable")

	config := InitConfiguration()
	assert.NotNil(t, config)
	assert.Equal(t, "localhost", config[DB_HOST])
	assert.Equal(t, "test_user", config[DB_USER])
}

func TestInitConfiguration_EnvNotSet(t *testing.T) {
	// Очистка всех переменных окружения
	for _, key := range []string{DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE} {
		os.Unsetenv(key)
	}

	config := InitConfiguration()
	assert.NotNil(t, config)
	assert.Empty(t, config[DB_HOST])
	assert.Empty(t, config[DB_USER])
}

func TestInitDB(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(DB_HOST, "localhost")
	t.Setenv(DB_USER, "postgres")
	t.Setenv(DB_PASSWORD, "5873")
	t.Setenv(DB_NAME, "visual_novel")
	t.Setenv(DB_PORT, "5432")
	t.Setenv(DB_SSLMODE, "prefer")

	db, err := InitDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestInitDB_InvalidDSN(t *testing.T) {
	// Установка неверного хоста
	t.Setenv(DB_HOST, "invalid_host")
	t.Setenv(DB_USER, "test_user")
	t.Setenv(DB_PASSWORD, "test_password")
	t.Setenv(DB_NAME, "test_db")
	t.Setenv(DB_PORT, "5432")
	t.Setenv(DB_SSLMODE, "disable")

	db, err := InitDB()
	assert.Error(t, err)
	assert.Nil(t, db)
}
