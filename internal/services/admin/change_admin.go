package admin

import (
	"gorm.io/gorm"
	"log"
	"vn/internal/storage"
)

const (
	NilStatus = -1
)

func ChangeAdmin(
	id int64,
	name string,
	email string,
	password string,
	adminStatus int,
	createdChapters []int64,
	db *gorm.DB,
) error {
	user, err := storage.SelectAdminWithId(db, id)

	log.Println(user)

	if err != nil {
		return err
	}

	if name != "" {
		user.Name = name
	}

	if email != "" {
		user.Email = email
	}

	if password != "" {
		user.Password = password
	}

	if adminStatus != NilStatus {
		user.AdminStatus = adminStatus
	}

	if createdChapters != nil {
		user.CreatedChapters = append(user.CreatedChapters, createdChapters...)
	}

	_, err = storage.UpdateAdmin(db, id, user)

	return err
}
