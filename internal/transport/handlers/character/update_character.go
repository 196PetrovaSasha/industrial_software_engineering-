package character

import (
	"db_novel_service/internal/services/character"
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

type UpdateCharacterRequest struct {
	Id       string            `json:"id"`
	Name     string            `json:"name,omitempty"`
	Slug     string            `json:"slug,omitempty"`
	Color    string            `json:"color,omitempty"`
	Emotions map[string]string `json:"emotions,omitempty"`
}

func UpdateCharacterHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("получен запрос на обновление персонажа 1")
		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		// Обрабатываем предварительный запрос (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		log.Println("получен запрос на обновление персонажа ")

		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req UpdateCharacterRequest
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("ошибка получения тела запроса на обноаление персонажа")
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Разбираем JSON
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println("ошибка распаковки запроса обновления персонажа")
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

		emotions := make(map[int64]int64)

		for i, j := range req.Emotions {
			emIndex, err := strconv.ParseInt(i, 10, 64)

			if err != nil {
				if err != nil {
					log.Println("ошибка конвертации index")
					http.Error(w, "Failed to covert index", http.StatusInternalServerError)
					return
				}
			}

			emottionId, err := strconv.ParseInt(j, 10, 64)

			if err != nil {
				if err != nil {
					log.Println("ошибка конвертации index")
					http.Error(w, "Failed to covert index", http.StatusInternalServerError)
					return
				}
			}

			emotions[emIndex] = emottionId
		}

		err = character.UpdateCharacter(id, req.Name, req.Slug, req.Color, emotions, db)

		log.Println(err)

		if err != nil {
			http.Error(w, "fail to create character", http.StatusInternalServerError)
		}
	}
}
