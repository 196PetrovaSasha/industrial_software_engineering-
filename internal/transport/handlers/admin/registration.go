package admin

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vn/internal/services/admin"
	"vn/pkg/metrick"
)

type AdminRegistrationRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func AdminRegistrationHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
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

		log.Info().Msg("получен запрос на регистрацию админа")
		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		// Обрабатываем предварительный запрос (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем метод запроса
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			log.Error().Msg("Only POST requests allowed in registration admin")
			return
		}

		// Остальной код обработки POST-запроса остается без изменений
		var req AdminRegistrationRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("Failed to read request body in registration admin")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in registration admin")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		id, err := admin.Registration(req.Email, req.Name, req.Password, db)
		if err != nil {
			log.Error().Msg("fail to register admin in registration admin")
			http.Error(w, "fail to register admin", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"id": id,
		}
		json.NewEncoder(w).Encode(response)
	}
}

type Response struct {
	id int64
}
