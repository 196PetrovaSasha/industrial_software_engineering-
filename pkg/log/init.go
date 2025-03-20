package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
)

// Environment variables and their default values
const (
	LOG_FILE     = "LOG_FILE"
	MAX_SIZE_MB  = "LOG_MAX_SIZE_MB"
	MAX_BACKUPS  = "LOG_MAX_BACKUPS"
	MAX_AGE_DAYS = "LOG_MAX_AGE_DAYS"
	COMPRESS     = "LOG_COMPRESS"
	LOG_LEVEL    = "LOG_LEVEL"
	DEBUG_MODE   = "DEBUG_MODE"
)

type LoggerConfig struct {
	Filename   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	Level      string
	DebugMode  bool
}

func NewLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Filename:   getEnv(LOG_FILE, "file.log"),
		MaxSizeMB:  getEnvInt(MAX_SIZE_MB, 100),
		MaxBackups: getEnvInt(MAX_BACKUPS, 5),
		MaxAgeDays: getEnvInt(MAX_AGE_DAYS, 30),
		Compress:   getEnvBool(COMPRESS, true),
		Level:      getEnv(LOG_LEVEL, "info"),
		DebugMode:  getEnvBool(DEBUG_MODE, false),
	}
}

func NewLogger() *zerolog.Logger {
	config := NewLoggerConfig()
	logLevel, _ := zerolog.ParseLevel(config.Level)

	logger := zerolog.New(zerolog.New(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   config.Compress,
	}))

	logger = logger.With().Timestamp().Logger()
	logger = logger.Level(logLevel)

	return &logger
}

func getEnv(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultVal
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, _ := strconv.Atoi(getEnv(key, fmt.Sprintf("%d", defaultValue)))
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value, _ := strconv.ParseBool(getEnv(key, fmt.Sprintf("%t", defaultValue)))
	return value
}
