package storage

import (
	"errors"
	"gorm.io/gorm"
	"vn/internal/models"
)

func RegisterMedia(db *gorm.DB, media models.Media) (int64, error) {
	result := db.Create(&media)
	if result.RowsAffected == 0 {
		return 0, errors.New("media not created")
	}
	return result.RowsAffected, nil
}

func SelectMediaWIthId(db *gorm.DB, id int64) (models.Media, error) {
	var media models.Media
	result := db.First(&media, "id = ?", id)
	if result.RowsAffected == 0 {
		return models.Media{}, errors.New("media data not found")
	}
	return media, nil
}

func SelectMedia(db *gorm.DB) ([]models.Media, error) {
	var media []models.Media
	result := db.Find(&media)

	if result.Error != nil {
		return nil, result.Error
	}

	return media, nil
}

func DeleteMedia(db *gorm.DB, id int64) (int64, error) {
	var deletedMedia models.Media
	result := db.Where("id = ?", id).Delete(&deletedMedia)
	if result.RowsAffected == 0 {
		return 0, errors.New("media data not update")
	}
	return result.RowsAffected, nil
}
