package chapter

import (
	"gorm.io/gorm"
	"log"
	"vn/internal/models"
	"vn/internal/storage"
)

func GetChaptersByUserId(db *gorm.DB, id int64) ([]models.Chapter, error) {

	log.Println(id)
	_, err := storage.SelectAdminWithId(db, id)

	log.Println(err)

	if err == nil {
		chapters, err := storage.GetChaptersForAdmin(db)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		return chapters, nil
	}

	_, err = storage.SelectPlayerWIthId(db, id)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	chapters, err := storage.FindPublishedChapters(db)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return chapters, nil
}
