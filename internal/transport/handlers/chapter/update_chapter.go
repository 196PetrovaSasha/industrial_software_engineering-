package chapter

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"vn/internal/services/chapter"
	"vn/pkg/metrick"
)

type UpdateChapterRequest struct {
	Id             string   `json:"id"`
	Name           string   `json:"name,omitempty"`
	StartNode      string   `json:"start_node,omitempty"`
	Nodes          []string `json:"nodes,omitempty"`
	Characters     []string `json:"characters,omitempty"`
	Status         int      `json:"status,omitempty"` // 0 - черновик, 1 - на проверке, 2 - опубликована
	UpdateAuthorId string   `json:"update_author_id,omitempty"`
}

func UpdateChapterHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
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

		log.Info().Msg("получен запрос на обновление главы")

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
			log.Error().Msg("Only POST requests allowed in chapters update")
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req UpdateChapterRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("Failed to read request body in chapters update")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Error().Msg("Invalid JSON format in chapters update")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(req.Id, 10, 64)

		if err != nil {
			if err != nil {
				log.Error().Msg("Failed to covert id in chapters update")
				http.Error(w, "Failed to covert id", http.StatusInternalServerError)
				return
			}
		}

		var nodes []int64

		for _, node := range req.Nodes {
			nodeId, err := strconv.ParseInt(node, 10, 64)

			if err != nil {
				if err != nil {
					log.Error().Msg("Failed to covert id in chapters update")
					http.Error(w, "Failed to covert id", http.StatusInternalServerError)
					return
				}
			}

			nodes = append(nodes, nodeId)
		}

		var characters []int64

		for _, character := range req.Characters {
			characterId, err := strconv.ParseInt(character, 10, 64)

			if err != nil {
				if err != nil {
					log.Error().Msg("Failed to covert id in chapters update")
					http.Error(w, "Failed to covert id", http.StatusInternalServerError)
					return
				}
			}

			characters = append(characters, characterId)
		}

		author, err := strconv.ParseInt(req.UpdateAuthorId, 10, 64)

		if err != nil {
			if err != nil {
				log.Error().Msg("Failed to covert id in chapters update")
				http.Error(w, "Failed to covert id", http.StatusInternalServerError)
				return
			}
		}

		log.Println("startNode", req.StartNode)

		var startNode int64

		if req.StartNode != "" {
			startNode, err = strconv.ParseInt(req.StartNode, 10, 64)

			if err != nil {
				if err != nil {
					log.Error().Msg("Failed to covert id in chapters update")
					http.Error(w, "Failed to covert id", http.StatusInternalServerError)
					return
				}
			}
		} else {
			startNode = 0
		}

		err = chapter.UpdateChapter(id, req.Name, nodes, characters, author, startNode, req.Status, db)

		if err != nil {
			log.Error().Msg("fail to create chapter in chapters update")
			http.Error(w, "fail to create chapter", http.StatusInternalServerError)
		}
	}
}
