package main

import (
	"context"
	"db_novel_service/cmd/service/model"
	"db_novel_service/internal/transport/handlers/admin"
	"db_novel_service/internal/transport/handlers/chapter"
	"db_novel_service/internal/transport/handlers/character"
	"db_novel_service/pkg/atlas"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	PORT = "APP_PORT"
)

func main() {

	service := run()

	correctDB := atlas.StartAtlasSchemaValidation()

	if !correctDB {
		service.Log.Error().Msg("ошибка валидации схемы базы данных")
		return
	} else {
		service.Log.Info().Msg("валидация схем баз данных прошла успешно")
	}

	generateSQLMetadata(service.DB) // наглядный вывод информации по бд

	service.Router.HandleFunc("/create-chapter", func(w http.ResponseWriter, r *http.Request) {
		handler := chapter.CreateChapterHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/update-chapter", func(w http.ResponseWriter, r *http.Request) {
		handler := chapter.UpdateChapterHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/get-chapters", func(w http.ResponseWriter, r *http.Request) {
		handler := chapter.GetChaptersByUserIdHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})

	service.Router.HandleFunc("/admin-authorization", func(w http.ResponseWriter, r *http.Request) {
		handler := admin.AdminAuthorisationHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/admin-registration", func(w http.ResponseWriter, r *http.Request) {
		handler := admin.AdminRegistrationHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/admin-change", func(w http.ResponseWriter, r *http.Request) {
		handler := admin.ChangeAdminHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})

	service.Router.HandleFunc("/create-character", func(w http.ResponseWriter, r *http.Request) {
		handler := character.CreateCharacterHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/update-character", func(w http.ResponseWriter, r *http.Request) {
		handler := character.UpdateCharacterHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})
	service.Router.HandleFunc("/get-characters", func(w http.ResponseWriter, r *http.Request) {
		handler := character.GetCharacterHandler(service.DB, service.Log)
		handler.ServeHTTP(w, r)
	})

	// Создаем экземпляр сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: service.Router,
	}

	// Регистрируем обработчик сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка сервера: %v", err)
		}
		log.Println("Сервер завершил обработку новых подключений")
	}()

	// Ждем сигнала завершения
	log.Println("Сервер запущен. Нажмите Ctrl+C для завершения...")
	<-stop

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("ошибка при завершении работы: %v", err)
	}

	log.Println("Сервер успешно завершен")
}

func run() *model.Service {
	service := model.NewService()

	service.Log.Info().Msg("service is created ")

	return service
}

func generateSQLMetadata(db *gorm.DB) error {
	// Получение схемы
	migrator := db.Migrator()
	if migrator == nil {
		log.Fatal("Ошибка: migrator не найден")
	}

	// Получение информации о таблицах
	tables, err := migrator.GetTables()
	if err != nil {
		log.Fatal(err)
	}

	// Вывод информации о таблицах
	for _, table := range tables {
		log.Printf("Таблица: %s", table)

		// Получение информации о колонках
		columns, err := migrator.ColumnTypes(table)
		if err != nil {
			log.Printf("Ошибка получения колонок для %s: %v", table, err)
			continue
		}

		for _, column := range columns {
			log.Printf("  - Колонка: %s (%s)", column.Name(), column.DatabaseTypeName())
		}
	}

	return nil
}
