package storage

import (
	"db_novel_service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func RegisterCharacter(db *gorm.DB, character models.Character) (int64, error) {
	result := db.Create(&character)
	if result.RowsAffected == 0 {
		return 0, errors.New("character not created")
	}
	return result.RowsAffected, nil
}

func SelectCharacterWIthId(db *gorm.DB, id int64) (models.Character, error) {
	// Используем raw SQL с явной обработкой JSON
	query := `
        SELECT id, name, slug, color, 
               CAST(emotions AS TEXT) as emotions_raw
        FROM characters
        WHERE id = ?
    `

	rows, err := db.Raw(query, id).Rows()
	if err != nil {
		return models.Character{}, err
	}
	defer rows.Close()

	var c models.Character
	var emotionsRaw string

	// Проверяем, есть ли данные
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return models.Character{}, fmt.Errorf("failed to read rows: %w", err)
		}
		return models.Character{}, gorm.ErrRecordNotFound
	}

	// Читаем данные только если они есть
	err = rows.Scan(&c.Id, &c.Name, &c.Slug, &c.Color, &emotionsRaw)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to scan row: %w", err)
	}

	// Десериализуем JSON в карту
	err = json.Unmarshal([]byte(emotionsRaw), &c.Emotions)
	if err != nil {
		return models.Character{}, fmt.Errorf("failed to unmarshal emotions: %w", err)
	}

	return c, rows.Err()
}

//func SelectCharacters(db *gorm.DB) ([]models.Character, error) {
//	var characters []models.Character
//	result := db.Find(&characters)
//
//	if result.Error != nil {
//		return nil, result.Error
//	}
//
//	return characters, nil
//}

func SelectCharacters(db *gorm.DB) ([]models.Character, error) {
	var characters []models.Character

	// Используем raw SQL с явной обработкой JSON
	query := `
        SELECT id, name, slug, color, 
               CAST(emotions AS TEXT) as emotions_raw
        FROM characters
    `

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var c models.Character
		var emotionsRaw string

		// Читаем все поля, включая сырые данные JSON
		err = rows.Scan(&c.Id, &c.Name, &c.Slug, &c.Color, &emotionsRaw)
		if err != nil {
			return nil, err
		}

		// Десериализуем JSON в карту
		err = json.Unmarshal([]byte(emotionsRaw), &c.Emotions)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal emotions: %w", err)
		}

		characters = append(characters, c)
	}

	return characters, rows.Err()
}

func UpdateCharacter(db *gorm.DB, id int64, newCharacter models.Character) (models.Character, error) {
	var character models.Character

	emotionsJSON, err := json.Marshal(newCharacter.Emotions)
	if err != nil {
		return models.Character{}, err
	}

	result := db.Model(&character).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"Name":     newCharacter.Name,
			"Slug":     newCharacter.Slug,
			"Color":    newCharacter.Color,
			"Emotions": json.RawMessage(emotionsJSON),
		})

	if result.RowsAffected == 0 {
		return models.Character{}, errors.New("character data not update")
	}

	return character, nil
}

func DeleteCharacter(db *gorm.DB, id int64) (int64, error) {
	var deletedCharacter models.Chapter
	result := db.Where("id = ?", id).Delete(&deletedCharacter)
	if result.RowsAffected == 0 {
		return 0, errors.New("character data not update")
	}
	return result.RowsAffected, nil
}
