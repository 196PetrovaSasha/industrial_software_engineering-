package migrator

import (
	"log"
	"vn/internal/models"
)

func AutoMigrate() {
	MigrateAdmin()
	MigratePlayer()
	MigrateChapter()
	MigrateCharacters()
	MigrateNode()
	MigrateMedia()
	MigrateRequest()
}

func MigrateAdmin() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Admin{})

	log.Println("Таблицы успешно созданы")
}

func MigratePlayer() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Player{})

	log.Println("Таблицы успешно созданы")
}

func MigrateChapter() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Chapter{})

	log.Println("Таблицы успешно созданы")
}

func MigrateCharacters() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Character{})

	log.Println("Таблицы успешно созданы")
}

func MigrateNode() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Node{}, &models.Branching{}, &models.EndInfo{})

	log.Println("Таблицы успешно созданы")
}

func MigrateMedia() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Media{})

	log.Println("Таблицы успешно созданы")
}

func MigrateRequest() {
	// Подключение к базе данных
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создание таблиц
	// При необходимрсти меняй на другой метод
	db.AutoMigrate(&models.Request{})

	log.Println("Таблицы успешно созданы")
}
