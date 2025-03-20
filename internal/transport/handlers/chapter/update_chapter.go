package chapter

import (
	"db_novel_service/internal/services/chapter"
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
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

		log.Println("получен запрос на изменение chapter")

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
		var req UpdateChapterRequest
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

		id, err := strconv.ParseInt(req.Id, 10, 64)

		if err != nil {
			if err != nil {
				log.Println("ошибка конвертации id")
				http.Error(w, "Failed to covert id", http.StatusInternalServerError)
				return
			}
		}

		var nodes []int64

		for _, node := range req.Nodes {
			nodeId, err := strconv.ParseInt(node, 10, 64)

			if err != nil {
				if err != nil {
					log.Println("ошибка конвертации nodeId")
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
					log.Println("ошибка конвертации chapterId")
					http.Error(w, "Failed to covert id", http.StatusInternalServerError)
					return
				}
			}

			characters = append(characters, characterId)
		}

		author, err := strconv.ParseInt(req.UpdateAuthorId, 10, 64)

		if err != nil {
			if err != nil {
				log.Println("ошибка конвертации authorId")
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
					log.Println("ошибка конвертации startNode")
					http.Error(w, "Failed to covert id", http.StatusInternalServerError)
					return
				}
			}
		} else {
			startNode = 0
		}

		err = chapter.UpdateChapter(id, req.Name, nodes, characters, author, startNode, req.Status, db)

		if err != nil {
			http.Error(w, "fail to create chapter", http.StatusInternalServerError)
		}
	}
}
