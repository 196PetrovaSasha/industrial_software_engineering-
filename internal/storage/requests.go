package storage

import (
	"db_novel_service/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func RegisterRequest(db *gorm.DB, request models.Request) (int64, error) {
	result := db.Create(&request)
	if result.RowsAffected == 0 {
		return 0, errors.New("node not created")
	}
	return result.RowsAffected, nil
}

func SelectRequestWithId(db *gorm.DB, id int64) (*models.Request, error) {
	var request models.Request
	err := db.First(&request, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("запрос с ID %d не найден", id)
	}
	return &request, err
}

func UpdateRequest(db *gorm.DB, id int64, newRequest models.Request) (int64, error) {
	result := db.Model(&newRequest).
		Where("id = ?", id).
		Omit("id"). // Исключаем ID из обновления
		Updates(newRequest)

	return result.RowsAffected, result.Error
}

func DeleteRequest(db *gorm.DB, id int64) (int64, error) {
	result := db.Where("id = ?", id).Delete(&models.Request{})
	return result.RowsAffected, result.Error
}

func GetAllRequests(db *gorm.DB) ([]models.Request, error) {
	var requests []models.Request

	// Используем Find вместо First для получения всех записей
	result := db.Find(&requests)

	// Проверяем ошибки
	if result.Error != nil {
		return []models.Request{}, fmt.Errorf("ошибка при получении запросов: %w", result.Error)
	}

	// Проверяем наличие записей
	if result.RowsAffected == 0 {
		return []models.Request{}, fmt.Errorf("записи не найдены")
	}

	return requests, nil
}
