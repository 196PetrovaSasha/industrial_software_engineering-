package config

import (
	"fmt"
	"os"
	"strconv"
)

// Environment variables and their default values
const (
	PORT_ENV     = "PORT"
	PORT_DEFAULT = ""
)

type Config struct {
	Port int64
}

func NewConfig() *Config {
	portStr := getEnv(PORT_ENV, PORT_DEFAULT)

	// Parse string to int64
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		return &Config{
			Port: 8080, // Default port if parsing fails
		}
	}

	return &Config{
		Port: port,
	}
}

// Helper function to get an environment variable or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	// Check if the key is a constant
	if _, ok := getEnvConstants()[key]; ok {
		return fmt.Sprintf("%v", getEnvConstants()[key])
	}

	return defaultVal
}

// Private function to map constants to their values
func getEnvConstants() map[string]interface{} {
	return map[string]interface{}{
		PORT_ENV:     "",
		PORT_DEFAULT: 8080,
	}
}
