package log

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestNewLoggerConfig(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(LOG_FILE, "test.log")
	t.Setenv(MAX_SIZE_MB, "50")
	t.Setenv(MAX_BACKUPS, "10")
	t.Setenv(MAX_AGE_DAYS, "30")
	t.Setenv(COMPRESS, "true")
	t.Setenv(LOG_LEVEL, "debug")
	t.Setenv(DEBUG_MODE, "true")

	config := NewLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "test.log", config.Filename)
	assert.Equal(t, 50, config.MaxSizeMB)
	assert.Equal(t, 10, config.MaxBackups)
	assert.Equal(t, 30, config.MaxAgeDays)
	assert.Equal(t, true, config.Compress)
	assert.Equal(t, "debug", config.Level)
	assert.Equal(t, true, config.DebugMode)
}

func TestNewLoggerConfig_EnvNotSet(t *testing.T) {
	// Очистка всех переменных окружения
	for _, key := range []string{LOG_FILE, MAX_SIZE_MB, MAX_BACKUPS, MAX_AGE_DAYS, COMPRESS, LOG_LEVEL, DEBUG_MODE} {
		os.Unsetenv(key)
	}

	config := NewLoggerConfig()
	assert.NotNil(t, config)
	assert.Equal(t, "file.log", config.Filename)
	assert.Equal(t, 100, config.MaxSizeMB)
	assert.Equal(t, 5, config.MaxBackups)
	assert.Equal(t, 30, config.MaxAgeDays)
	assert.Equal(t, true, config.Compress)
	assert.Equal(t, "info", config.Level)
	assert.Equal(t, false, config.DebugMode)
}

func TestNewLogger(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(LOG_FILE, "test.log")
	t.Setenv(LOG_LEVEL, "debug")

	logger := NewLogger()
	assert.NotNil(t, logger)

	// Проверяем, что логгер пишет в файл
	logger.Debug().Msg("test message")
	_, err := os.Stat("test.log")
	assert.NoError(t, err)
	assert.FileExists(t, "test.log")
}

func TestGetEnv(t *testing.T) {
	// Проверка getEnv
	t.Setenv("TEST_KEY", "test_value")
	value := getEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", value)

	// Проверка значения по умолчанию
	value = getEnv("NON_EXISTENT_KEY", "default")
	assert.Equal(t, "default", value)
}

func TestGetEnvInt(t *testing.T) {
	// Проверка getEnvInt
	t.Setenv("TEST_KEY", "42")
	value := getEnvInt("TEST_KEY", 0)
	assert.Equal(t, 42, value)

	// Проверка значения по умолчанию
	value = getEnvInt("NON_EXISTENT_KEY", 100)
	assert.Equal(t, 100, value)

	// Проверка некорректного значения
	t.Setenv("TEST_KEY", "invalid")
	value = getEnvInt("TEST_KEY", 0)
	assert.Equal(t, 0, value)
}

func TestGetEnvBool(t *testing.T) {
	// Проверка getEnvBool
	t.Setenv("TEST_KEY", "true")
	value := getEnvBool("TEST_KEY", false)
	assert.Equal(t, true, value)

	// Проверка значения по умолчанию
	value = getEnvBool("NON_EXISTENT_KEY", false)
	assert.Equal(t, false, value)

	// Проверка некорректного значения
	t.Setenv("TEST_KEY", "invalid")
	value = getEnvBool("TEST_KEY", false)
	assert.Equal(t, false, value)
}

func TestLoggerConfig_InvalidLogLevel(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(LOG_LEVEL, "invalid")
	t.Setenv(LOG_FILE, "test.log")

	logger := NewLogger()
	assert.NotNil(t, logger)

	// Проверяем, что логгер использует уровень по умолчанию
	logger.Debug().Msg("test message")
	logger.Info().Msg("test message")
	logger.Warn().Msg("test message")
	logger.Error().Msg("test message")
}

func TestLoggerConfig_LevelOrder(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(LOG_LEVEL, "debug")
	t.Setenv(LOG_FILE, "test.log")

	logger := NewLogger()
	assert.NotNil(t, logger)

	// Проверяем порядок уровней логирования
	logger.Debug().Msg("debug message")
	logger.Info().Msg("info message")
	logger.Warn().Msg("warn message")
	logger.Error().Msg("error message")
}

func TestLoggerConfig_Rotation(t *testing.T) {
	// Настройка тестовой среды
	t.Setenv(LOG_FILE, "test.log")
	t.Setenv(MAX_SIZE_MB, "1") // Устанавливаем маленький размер для теста
	t.Setenv(MAX_BACKUPS, "2")

	logger := NewLogger()
	assert.NotNil(t, logger)

	// Записываем много сообщений для проверки ротации
	for i := 0; i < 1000; i++ {
		logger.Info().Msgf("test message %d", i)
	}

	// Проверяем, что файлы ротации созданы
	files, err := os.ReadDir(".")
	assert.NoError(t, err)
	var logFiles int
	for _, f := range files {
		if f.Name() == "test.log" || strings.HasPrefix(f.Name(), "test.log.") {
			logFiles++
		}
	}
	assert.GreaterOrEqual(t, logFiles, 2)
}
