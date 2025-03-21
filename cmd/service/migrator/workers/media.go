package main

import (
	"github.com/joho/godotenv"
	"log"
	"vn/cmd/service/migrator"
	"vn/internal/models"
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
	db.AutoMigrate(&models.Media{})

	log.Println("Таблицы успешно созданы")
}
