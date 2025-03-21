package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"
	"vn/internal/models"
)

func RegisterPlayer(db *gorm.DB, player models.Player) (int64, error) {
	// Проверка на пустой email
	if strings.TrimSpace(player.Email) == "" {
		return 0, errors.New("email не может быть пустым")
	}

	// Создаем новый экземпляр структуры с инициализированными полями
	newPlayer := models.Player{
		Id:                player.Id,
		CompletedChapters: []int64{},
		ChaptersProgress:  map[int64]int64{},
		Email:             player.Email,
		Name:              player.Name,
		Phone:             player.Phone,
		Password:          player.Password,
		Admin:             player.Admin,
		SoundSettings:     player.SoundSettings,
	}

	// Создаем мапу значений для обновления
	updates := map[string]interface{}{
		"completed_chapters": json.RawMessage(`[]`),
		"chapters_progress":  json.RawMessage(`{}`),
	}

	result := db.Create(&newPlayer).
		Set("gorm:insert_option", "ON DUPLICATE KEY UPDATE").
		Updates(updates)

	if err := result.Error; err != nil {
		return 0, fmt.Errorf("ошибка при создании игрока: %w", err)
	}

	return result.RowsAffected, nil
}

//func SelectPlayerWIthEmail(db *gorm.DB, email string) (models.Player, error) {
//	var player models.Player
//	result := db.First(&player, "email = ?", email)
//	if result.RowsAffected == 0 {
//		return models.Player{}, errors.New("player data not found")
//	}
//	return player, nil
//}

func SelectPlayerWIthEmail(db *gorm.DB, email string) (models.Player, error) {
	var player models.Player

	// Используем raw SQL с явной обработкой JSON полей
	query := `
        SELECT id, name, email, phone, password, admin,
               CAST(completed_chapters AS TEXT) as completed_chapters_raw,
               CAST(chapters_progress AS TEXT) as chapters_progress_raw,
               sound_settings
        FROM players
        WHERE email = $1
        LIMIT 1
    `

	row := db.Raw(query, email).Row()
	if err := row.Err(); err != nil {
		return models.Player{}, err
	}

	var (
		completedChaptersRaw string
		chaptersProgressRaw  string
	)

	err := row.Scan(
		&player.Id,
		&player.Name,
		&player.Email,
		&player.Phone,
		&player.Password,
		&player.Admin,
		&completedChaptersRaw,
		&chaptersProgressRaw,
		&player.SoundSettings,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Player{}, errors.New("player data not found")
		}
		return models.Player{}, err
	}

	// Десериализуем JSON поля
	if completedChaptersRaw != "" {
		err = json.Unmarshal([]byte(completedChaptersRaw), &player.CompletedChapters)
		if err != nil {
			return models.Player{}, fmt.Errorf("failed to unmarshal completed_chapters: %w", err)
		}
	} else {
		player.CompletedChapters = []int64{}
	}

	if chaptersProgressRaw != "" {
		err = json.Unmarshal([]byte(chaptersProgressRaw), &player.ChaptersProgress)
		if err != nil {
			return models.Player{}, fmt.Errorf("failed to unmarshal chapters_progress: %w", err)
		}
	} else {
		player.ChaptersProgress = map[int64]int64{}
	}

	return player, nil
}

//func SelectPlayerWIthId(db *gorm.DB, id int64) (models.Player, error) {
//	var player models.Player
//	result := db.First(&player, "id = ?", id)
//	if result.RowsAffected == 0 {
//		return models.Player{}, errors.New("player data not found")
//	}
//	return player, nil
//}

func SelectPlayerWIthId(db *gorm.DB, id int64) (models.Player, error) {
	var player models.Player

	// Используем raw SQL с явной обработкой JSON полей
	query := `
        SELECT id, name, email, phone, password, admin,
               CAST(completed_chapters AS TEXT) as completed_chapters_raw,
               CAST(chapters_progress AS TEXT) as chapters_progress_raw,
               sound_settings
        FROM players
        WHERE id = $1
        LIMIT 1
    `

	row := db.Raw(query, id).Row()
	if err := row.Err(); err != nil {
		return models.Player{}, err
	}

	var (
		completedChaptersRaw string
		chaptersProgressRaw  string
	)

	err := row.Scan(
		&player.Id,
		&player.Name,
		&player.Email,
		&player.Phone,
		&player.Password,
		&player.Admin,
		&completedChaptersRaw,
		&chaptersProgressRaw,
		&player.SoundSettings,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Player{}, errors.New("player data not found")
		}
		return models.Player{}, err
	}

	// Десериализуем JSON поля
	if completedChaptersRaw != "" {
		err = json.Unmarshal([]byte(completedChaptersRaw), &player.CompletedChapters)
		if err != nil {
			return models.Player{}, fmt.Errorf("failed to unmarshal completed_chapters: %w", err)
		}
	} else {
		player.CompletedChapters = []int64{}
	}

	if chaptersProgressRaw != "" {
		err = json.Unmarshal([]byte(chaptersProgressRaw), &player.ChaptersProgress)
		if err != nil {
			return models.Player{}, fmt.Errorf("failed to unmarshal chapters_progress: %w", err)
		}
	} else {
		player.ChaptersProgress = map[int64]int64{}
	}

	return player, nil
}

func UpdatePlayer(db *gorm.DB, id int64, newPlayer models.Player) (models.Player, error) {
	var player models.Player

	// Логирование входящих данных для отладки
	log.Printf("UpdatePlayer: входящие данные:\nCompletedChapters: %+v\nChaptersProgress: %+v",
		newPlayer.CompletedChapters, newPlayer.ChaptersProgress)

	// Преобразуем массив пройденных глав в JSON
	completedChaptersJSON, err := json.Marshal(newPlayer.CompletedChapters)
	if err != nil {
		return models.Player{}, fmt.Errorf("ошибка маршалинга CompletedChapters: %w", err)
	}

	// Логирование JSON после маршалинга
	log.Printf("Маршаленный JSON для CompletedChapters: %s", completedChaptersJSON)

	// Преобразуем мапу прогресса в JSON
	chaptersProgressJSON, err := json.Marshal(newPlayer.ChaptersProgress)
	if err != nil {
		return models.Player{}, fmt.Errorf("ошибка маршалинга ChaptersProgress: %w", err)
	}

	// Обновляем данные игрока
	result := db.Model(&player).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":               newPlayer.Name,
			"email":              newPlayer.Email,
			"phone":              newPlayer.Phone,
			"password":           newPlayer.Password,
			"admin":              newPlayer.Admin,
			"completed_chapters": json.RawMessage(completedChaptersJSON),
			"chapters_progress":  json.RawMessage(chaptersProgressJSON),
			"sound_settings":     newPlayer.SoundSettings,
		})

	if result.RowsAffected == 0 {
		return models.Player{}, errors.New("player data not updated")
	}

	return player, nil
}

// Для отладки можно также добавить функцию проверки формата JSON
func validateJSONFormat(jsonData []byte, fieldName string) error {
	var raw map[string]interface{}
	err := json.Unmarshal(jsonData, &raw)
	if err != nil {
		return fmt.Errorf("неверный формат JSON для поля %s: %w", fieldName, err)
	}
	return nil
}

func DeletePlayer(db *gorm.DB, email string) (int64, error) {
	var deletedPlayer models.Node
	result := db.Where("email = ?", email).Delete(&deletedPlayer)
	if result.RowsAffected == 0 {
		return 0, errors.New("player data not update")
	}
	return result.RowsAffected, nil
}
