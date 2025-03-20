package character

import (
	"db_novel_service/internal/storage"
	"gorm.io/gorm"
	"log"
)

func UpdateCharacter(
	id int64,
	name string,
	slug string,
	color string,
	emotions map[int64]int64,
	db *gorm.DB,
) error {
	character, err := storage.SelectCharacterWIthId(db, id)

	log.Println(character)

	if err != nil {
		return err
	}

	newCharacter := character

	if name != "" {
		newCharacter.Name = name
	}

	if slug != "" {
		newCharacter.Slug = slug
	}

	if color != "" {
		newCharacter.Color = color
	}

	if emotions != nil {
		newCharacter.Emotions = emotions
	}

	_, err = storage.UpdateCharacter(db, id, newCharacter)

	log.Println(err)

	return err
}
