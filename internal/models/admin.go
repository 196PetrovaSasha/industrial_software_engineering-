package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Admin struct {
	Id               int64 `gorm:"primary_key"`
	Name             string
	Email            string
	Password         string
	AdminStatus      int     // 0 - дефолтный алмин, 1 - сверхадмин, -1 - незаапрувенный админ
	CreatedChapters  []int64 `gorm:"type:json;column:created_chapters"`
	RequestSent      []int64 `gorm:"type:json;column:request_sent"`
	RequestsReceived []int64 `gorm:"type:json;column:requests_received"`
}

func (admin *Admin) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	// Конвертируем []uint8 в json.RawMessage для JSON полей
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("неподдерживаемый тип данных при сканировании")
	}

	var data map[string]interface{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return fmt.Errorf("ошибка распаковки JSON: %w", err)
	}

	// Обновляем поля структуры из JSON
	admin.Name = data["name"].(string)
	admin.Email = data["email"].(string)
	admin.Password = data["password"].(string)
	admin.AdminStatus = int(data["admin_status"].(float64))

	// Конвертируем массивы из JSON
	createdChaptersRaw, ok := data["created_chapters"].([]interface{})
	if ok {
		admin.CreatedChapters = make([]int64, len(createdChaptersRaw))
		for i, v := range createdChaptersRaw {
			admin.CreatedChapters[i] = int64(v.(float64))
		}
	}

	requestSentRaw, ok := data["request_sent"].([]interface{})
	if ok {
		admin.RequestSent = make([]int64, len(requestSentRaw))
		for i, v := range requestSentRaw {
			admin.RequestSent[i] = int64(v.(float64))
		}
	}

	requestsReceivedRaw, ok := data["requests_received"].([]interface{})
	if ok {
		admin.RequestsReceived = make([]int64, len(requestsReceivedRaw))
		for i, v := range requestsReceivedRaw {
			admin.RequestsReceived[i] = int64(v.(float64))
		}
	}

	return nil
}

func (admin Admin) Value() (driver.Value, error) {
	data := map[string]interface{}{
		"name":              admin.Name,
		"email":             admin.Email,
		"password":          admin.Password,
		"admin_status":      admin.AdminStatus,
		"created_chapters":  admin.CreatedChapters,
		"request_sent":      admin.RequestSent,
		"requests_received": admin.RequestsReceived,
	}

	return json.Marshal(data)
}
