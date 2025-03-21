package character

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vn/internal/services/character"
	"vn/pkg/metrick"
)

type CreateCharacterRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func CreateCharacterHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		// Создаем wrapper для ResponseWriter чтобы отслеживать статус
		rw := &metrick.StatusRecorder{ResponseWriter: w}

		// Вызываем оригинальную функцию обработчика
		defer func() {
			// Записываем время выполнения запроса
			duration := time.Since(startTime).Seconds()

			// Записываем метрики
			metrick.RequestDuration.WithLabelValues("character", r.Method).
				Observe(duration)

			metrick.RequestCount.WithLabelValues(
				"character",
				r.Method,
				strconv.Itoa(rw.StatusCode),
			).Inc()
		}()

		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		// Обрабатываем предварительный запрос (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req CreateCharacterRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		id, err := character.CreateCharacter(req.Name, req.Slug, db)

		if err != nil {
			http.Error(w, "fail to create character", http.StatusInternalServerError)
		}

		// Формируем ответ
		response := map[string]interface{}{
			"id": utils.ToString(id),
		}

		log.Println("новый персонаж", id)

		// Отправляем ответ клиенту
		json.NewEncoder(w).Encode(response)
	}
}
