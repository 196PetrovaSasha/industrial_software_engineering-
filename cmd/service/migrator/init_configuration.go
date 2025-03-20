package migrator

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// Environment variables
const (
	DB_HOST     = "DB_HOST"
	DB_USER     = "DB_USER"
	DB_PASSWORD = "DB_PASSWORD"
	DB_NAME     = "DB_NAME"
	DB_PORT     = "DB_PORT"
	DB_SSLMODE  = "DB_SSLMODE"
)

func InitConfiguration() map[string]string {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	return map[string]string{
		DB_HOST:     os.Getenv(DB_HOST),
		DB_USER:     os.Getenv(DB_USER),
		DB_PASSWORD: os.Getenv(DB_PASSWORD),
		DB_NAME:     os.Getenv(DB_NAME),
		DB_PORT:     os.Getenv(DB_PORT),
		DB_SSLMODE:  os.Getenv(DB_SSLMODE),
	}
}

func InitDB() (*gorm.DB, error) {
	config := InitConfiguration()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config[DB_HOST], config[DB_USER], config[DB_PASSWORD],
		config[DB_NAME], config[DB_PORT], config[DB_SSLMODE])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
