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

type ChangeAdminRequest struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name,omitempty"`
	Email           string  `json:"email,omitempty"`
	Password        string  `json:"password,omitempty"`
	AdminStatus     int     `json:"status,omitempty"`
	CreatedChapters []int64 `json:"created_chapters,omitempty"`
}

func ChangeAdminHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
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

		log.Info().Msg("получен запрос на измение админа")
		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req ChangeAdminRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("Failed to read request body in update admin")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in update admin")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Здесь должна быть логика получения данных пользователя
		// Например, из базы данных:
		err = admin.ChangeAdmin(req.Id, req.Name, req.Email, req.Password, req.AdminStatus, req.CreatedChapters, db)

		if err != nil {
			log.Error().Msg("fail to change admin in update admin")
			http.Error(w, "fail to change admin", http.StatusInternalServerError)
		}

		log.Print(err)

		// Формируем ответ
		response := map[string]interface{}{
			"id": req.Id,
		}

		// Отправляем ответ клиенту
		json.NewEncoder(w).Encode(response)
	}
}
