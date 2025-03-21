package character

import (
	"gorm.io/gorm"
	"vn/internal/models"
	"vn/internal/storage"
)

func GetCharacters(db *gorm.DB) (*[]models.Character, error) {

	characters, err := storage.SelectCharacters(db)

	if err != nil {
		return nil, err
	}

	return &characters, nil
}
