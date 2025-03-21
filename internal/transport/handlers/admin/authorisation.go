package admin

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vn/internal/services/admin"
	"vn/pkg/metrick"
)

type UserAuthorisationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminAuthorisationHandler(db *gorm.DB, log *zerolog.Logger, authConfig AuthConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		// Создаем wrapper для ResponseWriter чтобы отслеживать статус
		rw := &metrick.StatusRecorder{ResponseWriter: w}

		// Вызываем оригинальную функцию обработчика
		defer func() {
			// Записываем время выполнения запроса
			duration := time.Since(startTime).Seconds()

			// Записываем метрики
			metrick.RequestDuration.WithLabelValues("admin", r.Method).
				Observe(duration)

			metrick.RequestCount.WithLabelValues(
				"admin",
				r.Method,
				strconv.Itoa(rw.StatusCode),
			).Inc()
		}()

		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")

		// Обрабатываем предварительный запрос (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			log.Error().Msg("invalid request type in authorization admin")
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UserAuthorisationRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("failed to read request body in authorization admin")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in authorization admin")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			log.Error().Msg("Email and password are required in authorization admin")
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		// Получаем данные пользователя
		user, err := admin.Authorization(req.Email, req.Password, db)
		if err != nil {
			log.Error().Msg("Authorization failed in authorization admin")
			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}

		if user.AdminStatus == -1 {
			http.Error(w, "Authorization failed", http.StatusForbidden)
			return
		}

		// Создаём JWT токен
		token, err := GenerateToken(utils.ToString(user.Id), "admin", authConfig)
		if err != nil {
			log.Error().Msg("Error generating JWT token")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Формируем ответ с токеном
		response := map[string]interface{}{
			"id":               utils.ToString(user.Id),
			"name":             user.Name,
			"adminStatus":      user.AdminStatus,
			"createdChapters":  user.CreatedChapters,
			"requestSent":      user.RequestSent,
			"requestsReceived": user.RequestsReceived,
			"token":            token,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

// GenerateToken генерирует новый JWT токен
func GenerateToken(userID string, role string, authConfig AuthConfig) (string, error) {
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

// AuthConfig содержит настройки аутентификации
type AuthConfig struct {
	SecretKey string
	TTL       time.Duration
}
