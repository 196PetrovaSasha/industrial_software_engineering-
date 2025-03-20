package admin

import (
	"db_novel_service/internal/services/admin"
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
)

type UserAuthorisationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminAuthorisationHandler(db *gorm.DB, log *zerolog.Logger) http.HandlerFunc {
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

		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close() // Важно закрыть тело запроса

		// Разбираем JSON
		var req UserAuthorisationRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println(err)

			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Проверяем, что поля не пустые
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		log.Println(req.Email, req.Password)

		// Получаем данные пользователя
		user, err := admin.Authorization(req.Email, req.Password, db)
		if err != nil {
			log.Println(err)

			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}

		if user.AdminStatus == -1 {
			http.Error(w, "Authorization failed", http.StatusForbidden)
		}

		log.Println(err)

		// Формируем ответ
		response := map[string]interface{}{
			"id":               utils.ToString(user.Id),
			"name":             user.Name,
			"adminStatus":      user.AdminStatus,
			"createdChapters":  user.CreatedChapters,
			"requestSent":      user.RequestSent,
			"requestsReceived": user.RequestsReceived,
		}

		log.Println("все прошло упешно", "авторизация", user.AdminStatus)

		// Отправляем ответ клиенту
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
