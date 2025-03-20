package main

import (
	"db_novel_service/cmd/service/migrator"
	"db_novel_service/internal/models"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Подключение к базе данных
	db, err := migrator.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Admin{})

	log.Println("Таблицы успешно созданы")
}
