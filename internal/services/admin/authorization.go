package admin

import (
	"db_novel_service/internal/models"
	"db_novel_service/internal/storage"
	"fmt"
	"gorm.io/gorm"
)

func Authorization(email string, password string, db *gorm.DB) (*models.Admin, error) {
	user, err := storage.SelectAdminWIthEmail(db, email)

	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
