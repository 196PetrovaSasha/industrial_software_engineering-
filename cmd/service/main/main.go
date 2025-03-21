package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"vn/cmd/service/model"
	"vn/internal/transport/handlers/admin"
	"vn/internal/transport/handlers/chapter"
	"vn/internal/transport/handlers/character"
	"vn/pkg/atlas"
	"vn/pkg/metrick"
)

var registry = prometheus.NewRegistry()

// init is invoked before main()
func init() {
	if err := registry.Register(metrick.RequestDuration); err != nil {
		log.Printf("Не удалось зарегистрировать RequestDuration: %v", err)
	}

	if err := registry.Register(metrick.RequestCount); err != nil {
		log.Printf("Не удалось зарегистрировать RequestCount: %v", err)
	}

	// Добавляем endpoint для метрик
	http.Handle("/metrics", promhttp.Handler())

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	service := run()

	correctDB := atlas.StartAtlasSchemaValidation()

	if !correctDB {
		service.Log.Error().Msg("ошибка валидации схемы базы данных")
		return
	} else {
		service.Log.Info().Msg("валидация схем баз данных прошла успешно")
	}

	generateSQLMetadata(service.DB, service.Log) // наглядный вывод информации по бд

	authConfig := admin.AuthConfig{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

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
		handler := admin.AdminAuthorisationHandler(service.DB, service.Log, authConfig)
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

	service.Log.Info().Msg("сервер успешно создан")

	// Регистрируем обработчик сигналов
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	service.Log.Info().Msg("обработчк сигналов успешно зарегестрирован")

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			service.Log.Fatal()
		}
		service.Log.Println("Сервер завершил обработку новых подключений")
	}()

	// Ждем сигнала завершения
	service.Log.Println("Сервер запущен. Нажмите Ctrl+C для завершения...")
	<-stop

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("ошибка при завершении работы: %v", err)
	}

	service.Log.Println("Сервер успешно завершен")
}

func run() *model.Service {
	service := model.NewService()

	// Загружаем конфигурацию из .env
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	ttlHours, _ := strconv.Atoi(os.Getenv("JWT_TTL_HOURS"))

	authConfig := admin.AuthConfig{
		SecretKey: jwtSecret,
		TTL:       time.Hour * time.Duration(ttlHours),
	}

	// Применяем middleware ко всем маршрутам
	service.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authMiddleware := AuthMiddleware(authConfig)
			authMiddleware(next.ServeHTTP)(w, r)
		})
	})

	return service
}

func generateSQLMetadata(db *gorm.DB, log *zerolog.Logger) error {
	// Получение схемы
	migrator := db.Migrator()
	if migrator == nil {
		log.Fatal().Msg("Ошибка: migrator не найден")
	}

	// Получение информации о таблицах
	tables, err := migrator.GetTables()
	if err != nil {
		log.Fatal().Msg(err.Error())
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

// AuthMiddleware проверяет JWT токен в каждом запросе
func AuthMiddleware(authConfig admin.AuthConfig) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем запросы к конечным точкам аутентификации
			if r.URL.Path == "/admin-authorization" || r.URL.Path == "/admin-registration" {
				next(w, r)
				return
			}

			// Получаем токен из заголовка
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				http.Error(w, "Токен аутентификации не предоставлен", http.StatusUnauthorized)
				return
			}

			// Проверяем токен
			token, err := verifyToken(tokenStr, authConfig.SecretKey)
			if err != nil {
				http.Error(w, "Неверный токен аутентификации", http.StatusUnauthorized)
				return
			}

			// Добавляем информацию о пользователе в контекст
			ctx := context.WithValue(r.Context(), "user", token.Claims)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

// verifyToken проверяет валидность JWT токена
func verifyToken(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка валидации токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("невалидный токен")
	}

	return token, nil
}

// GenerateToken генерирует новый JWT токен
func GenerateToken(userID string, role string, authConfig admin.AuthConfig) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(authConfig.TTL).Unix(),
	})

	tokenString, err := claims.SignedString([]byte(authConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании токена: %w", err)
	}

	return tokenString, nil
}

// AdminAuthorisationHandler обновленный обработчик авторизации
func AdminAuthorisationHandler(db *gorm.DB, log *zerolog.Logger, authConfig admin.AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		// Здесь должна быть логика проверки учетных данных
		// Для примера используем фиксированные данные
		if credentials.Username == "admin" && credentials.Password == "password" {
			token, err := GenerateToken(credentials.Username, "admin", authConfig)
			if err != nil {
				http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			return
		}

		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
	}
}
