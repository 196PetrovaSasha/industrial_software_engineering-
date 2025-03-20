package admin

import (
	"db_novel_service/internal/services/admin"
	"encoding/json"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
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
		// Проверяем, что это POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
			return
		}

		// Читаем тело запроса
		var req ChangeAdminRequest
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

		// Здесь должна быть логика получения данных пользователя
		// Например, из базы данных:
		err = admin.ChangeAdmin(req.Id, req.Name, req.Email, req.Password, req.AdminStatus, req.CreatedChapters, db)

		if err != nil {
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
