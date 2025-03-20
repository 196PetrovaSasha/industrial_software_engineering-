package storage

import (
	"database/sql"
	"db_novel_service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	gorm "gorm.io/gorm"
	"log"
)

func RegisterAdmin(db *gorm.DB, admin models.Admin) (int64, error) {
	result := db.Create(&admin)
	if result.RowsAffected == 0 {
		return 0, errors.New("admin not created")
	}
	return result.RowsAffected, nil
}

func SelectAdminWIthEmail(db *gorm.DB, email string) (models.Admin, error) {
	var admin models.Admin
	// Используем raw SQL с явной обработкой JSON
	query := `
        SELECT 
            id,
            name,
            email,
            password,
            admin_status,
            COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw,
            COALESCE(request_sent::TEXT, '[]') as request_sent_raw,
            COALESCE(requests_received::TEXT, '[]') as requests_received_raw
        FROM admins
        WHERE email = $1
        LIMIT 1
    `
	row := db.Raw(query, email).Row()
	if err := row.Err(); err != nil {
		return models.Admin{}, err
	}

	// Декларируем переменные для JSON полей
	var (
		createdChaptersRaw  sql.NullString
		requestSentRaw      sql.NullString
		requestsReceivedRaw sql.NullString
	)

	// Сканируем все поля
	err := row.Scan(
		&admin.Id,
		&admin.Name,
		&admin.Email,
		&admin.Password,
		&admin.AdminStatus,
		&createdChaptersRaw,
		&requestSentRaw,
		&requestsReceivedRaw,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, errors.New("admin data not found")
		}
		return models.Admin{}, err
	}

	if createdChaptersRaw.Valid {
		err = json.Unmarshal([]byte(createdChaptersRaw.String), &admin.CreatedChapters)
		if err != nil {
			return models.Admin{}, fmt.Errorf("failed to unmarshal created_chapters: %w", err)
		}
	} else {
		admin.CreatedChapters = []int64{}
	}

	if requestSentRaw.Valid {
		err = json.Unmarshal([]byte(requestSentRaw.String), &admin.RequestSent)
		if err != nil {
			return models.Admin{}, fmt.Errorf("failed to unmarshal request_sent: %w", err)
		}
	} else {
		admin.RequestSent = []int64{}
	}

	if requestsReceivedRaw.Valid {
		err = json.Unmarshal([]byte(requestsReceivedRaw.String), &admin.RequestsReceived)
		if err != nil {
			return models.Admin{}, fmt.Errorf("failed to unmarshal requests_received: %w", err)
		}
	} else {
		admin.RequestsReceived = []int64{}
	}

	return admin, nil
}

func SelectAdminWithId(db *gorm.DB, id int64) (models.Admin, error) {
	var admin models.Admin

	log.Println(id)

	// Используем raw SQL с явной обработкой всех JSON полей
	query := `
        SELECT id, name, email, password, admin_status,
               COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw,
               COALESCE(request_sent::TEXT, '[]') as request_sent_raw,
               COALESCE(requests_received::TEXT, '[]') as requests_received_raw
        FROM admins
        WHERE id = $1
        LIMIT 1
    `

	row := db.Raw(query, id).Row()
	if err := row.Err(); err != nil {
		return models.Admin{}, err
	}

	var (
		nameRaw             string
		emailRaw            string
		passwordRaw         string
		adminStatusRaw      int
		createdChaptersRaw  string
		requestSentRaw      string
		requestsReceivedRaw string
	)

	err := row.Scan(
		&admin.Id,
		&nameRaw,
		&emailRaw,
		&passwordRaw,
		&adminStatusRaw,
		&createdChaptersRaw,
		&requestSentRaw,
		&requestsReceivedRaw,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Admin{}, errors.New("admin data not found")
		}
		return models.Admin{}, err
	}

	// Преобразуем строки в структуру
	admin.Name = nameRaw
	admin.Email = emailRaw
	admin.Password = passwordRaw
	admin.AdminStatus = adminStatusRaw

	// Десериализуем JSON поля
	err = json.Unmarshal([]byte(createdChaptersRaw), &admin.CreatedChapters)
	if err != nil {
		return models.Admin{}, fmt.Errorf("failed to unmarshal created_chapters: %w", err)
	}

	err = json.Unmarshal([]byte(requestSentRaw), &admin.RequestSent)
	if err != nil {
		return models.Admin{}, fmt.Errorf("failed to unmarshal request_sent: %w", err)
	}

	err = json.Unmarshal([]byte(requestsReceivedRaw), &admin.RequestsReceived)
	if err != nil {
		return models.Admin{}, fmt.Errorf("failed to unmarshal requests_received: %w", err)
	}

	return admin, nil
}

func SelectAllSupeAdmins(db *gorm.DB) ([]models.Admin, error) {
	var admins []models.Admin
	result := db.Where("admin_status = ?", 1).Find(&admins)

	if result.RowsAffected == 0 {
		return nil, errors.New("no super admin found")
	}

	return admins, nil
}

func UpdateAdmin(db *gorm.DB, id int64, newAdmin models.Admin) (models.Admin, error) {
	var admin models.Admin

	chaptersJSON, err := json.Marshal(newAdmin.CreatedChapters)

	if err != nil {
		return models.Admin{}, err
	}

	requestsSendJSON, err := json.Marshal(newAdmin.RequestSent)

	if err != nil {
		return models.Admin{}, err
	}

	requestsReceivedJSON, err := json.Marshal(newAdmin.RequestsReceived)

	if err != nil {
		return models.Admin{}, err
	}

	result := db.Model(&admin).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email":             newAdmin.Email,
			"password":          newAdmin.Password,
			"created_chapters":  json.RawMessage(chaptersJSON),
			"admin_status":      newAdmin.AdminStatus,
			"name":              newAdmin.Name,
			"request_sent":      json.RawMessage(requestsSendJSON),
			"requests_received": json.RawMessage(requestsReceivedJSON),
		})

	if result.RowsAffected == 0 {
		return models.Admin{}, errors.New("admin data not update")
	}

	return admin, nil
}
