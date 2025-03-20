package character

import (
	"db_novel_service/internal/models"
	"db_novel_service/internal/storage"
	"gorm.io/gorm"
)

func GetCharacters(db *gorm.DB) (*[]models.Character, error) {

	characters, err := storage.SelectCharacters(db)

	if err != nil {
		return nil, err
	}

	return &characters, nil
}
