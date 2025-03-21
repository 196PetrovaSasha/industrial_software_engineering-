package chapter

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vn/internal/services/chapter"
	"vn/pkg/metrick"
)

type CreateChapterRequest struct {
	Author string `json:"author"`
}

func CreateChapterHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		// Создаем wrapper для ResponseWriter чтобы отслеживать статус
		rw := &metrick.StatusRecorder{ResponseWriter: w}

		// Вызываем оригинальную функцию обработчика
		defer func() {
			// Записываем время выполнения запроса
			duration := time.Since(startTime).Seconds()

			// Записываем метрики
			metrick.RequestDuration.WithLabelValues("chapter", r.Method).
				Observe(duration)

			metrick.RequestCount.WithLabelValues(
				"chapter",
				r.Method,
				strconv.Itoa(rw.StatusCode),
			).Inc()
		}()

		log.Info().Msg("получен запрос на создание главы")
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
			log.Error().Msg("Only POST requests allowed in create chapter")
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req CreateChapterRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("Failed to read request body in create chapter")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in create chapter")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(req.Author, 10, 64)

		if err != nil {
			if err != nil {
				log.Error().Msg("Failed to covert id in create chapter")
				http.Error(w, "Failed to covert id", http.StatusInternalServerError)
				return
			}
		}

		id, nodeId, err := chapter.CreateDefaultChapter(id, db)

		if err != nil {
			log.Error().Msg("fail to create chapter in create chapter")
			http.Error(w, "fail to create chapter", http.StatusInternalServerError)
		}

		// Формируем ответ
		response := map[string]interface{}{
			"id":         utils.ToString(id),
			"start_node": utils.ToString(nodeId),
		}

		// Отправляем ответ клиенту
		json.NewEncoder(w).Encode(response)
	}
}

type CreateChapterResponse struct {
	Id        string
	startNode string
}
