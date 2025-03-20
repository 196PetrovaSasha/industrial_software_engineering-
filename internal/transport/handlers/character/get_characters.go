package character

import (
	"db_novel_service/internal/models"
	"db_novel_service/internal/services/character"
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"net/http"
)

func GetCharacterHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Добавляем CORS заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		// Обрабатываем предварительный запрос (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем, что это GET-запрос
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
			return
		}

		characters, err := character.GetCharacters(db)

		log.Println(err)

		if err != nil {
			http.Error(w, "fail to get characters", http.StatusInternalServerError)
		}

		// Формируем ответ
		response := map[string]interface{}{
			"characters": PrepareCharacterForResponse(characters),
		}

		// Отправляем ответ клиенту
		json.NewEncoder(w).Encode(response)
	}
}

type ResponseCharacter struct {
	Id       string
	Name     string
	Slug     string
	Color    string
	Emotions map[string]string
}

func PrepareCharacterForResponse(characters *[]models.Character) []ResponseCharacter {
	var res []ResponseCharacter

	if characters == nil {
		return res
	}

	for _, char := range *characters {
		emotions := make(map[string]string, len(char.Emotions))

		for i, j := range char.Emotions {
			emotions[utils.ToString(i)] = utils.ToString(j)
		}

		res = append(res, ResponseCharacter{
			Id:       utils.ToString(char.Id),
			Name:     char.Name,
			Slug:     char.Slug,
			Color:    char.Color,
			Emotions: emotions,
		})
	}

	return res
}
