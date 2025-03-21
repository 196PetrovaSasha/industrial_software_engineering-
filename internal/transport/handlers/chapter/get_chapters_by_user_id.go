package chapter

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"vn/internal/models"
	"vn/internal/services/chapter"
	"vn/pkg/metrick"
)

type GetChaptersByUserIdRequest struct {
	UserId string `json:"user_id"`
}

func GetChaptersByUserIdHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
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

		log.Info().Msg("получен запрос на получение главы")
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
			log.Error().Msg("Only POST requests allowed шin chapters")
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req GetChaptersByUserIdRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("Failed to read request body in chapters")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in chapters")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(req.UserId, 10, 64)

		if err != nil {
			if err != nil {
				log.Error().Msg("Failed to covert id in chapters")
				http.Error(w, "Failed to covert id", http.StatusInternalServerError)
				return
			}
		}

		chapters, err := chapter.GetChaptersByUserId(db, id)

		if err != nil {
			log.Error().Msg("fail to get chapters in chapters")
			http.Error(w, "fail to get chapters", http.StatusInternalServerError)
			return // Добавлен return
		}

		log.Println(err)

		// Формируем ответ
		response := map[string]interface{}{
			"chapters": prepareChaptersForResponce(chapters),
		}

		// Отправляем ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type ResponceChapter struct {
	Id         string
	Name       string
	StartNode  string
	Nodes      []string
	Characters []string
	Status     int
	Author     string
}

func prepareChaptersForResponce(chapters []models.Chapter) []ResponceChapter {
	var res []ResponceChapter

	for _, ch := range chapters {

		var nodes []string
		var characters []string

		for _, n := range ch.Nodes {
			nodes = append(nodes, utils.ToString(n))
		}

		for _, char := range ch.Characters {
			characters = append(characters, utils.ToString(char))
		}

		res = append(res, ResponceChapter{
			Id:         utils.ToString(ch.Id),
			Name:       ch.Name,
			StartNode:  utils.ToString(ch.StartNode),
			Nodes:      nodes,
			Characters: characters,
			Status:     ch.Status,
			Author:     utils.ToString(ch.Author),
		})
	}

	log.Println(len(res))

	return res
}

type ResponseChapter struct {
	Id         string `json:"id"`
	Name       string
	StartNode  string
	Nodes      []string
	Characters []string
	Status     int
	Author     string
}
