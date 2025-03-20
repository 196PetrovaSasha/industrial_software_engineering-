package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Конфигурация подключения к базе данных
	pgUser := "postgres"
	pgPassword := "5873"
	pgHost := "127.0.0.1"
	pgPort := "5432"
	pgDBName := "visual_novel"

	// Подключение к базе данных
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	// Получаем список таблиц
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
	if err != nil {
		log.Fatalf("Ошибка при выполнении запроса: %v", err)
	}

	defer rows.Close()

	fmt.Println("\nСписок таблиц:")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v", err)
		}
		fmt.Println(tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Ошибка при обходе строк: %v", err)
	}
}
