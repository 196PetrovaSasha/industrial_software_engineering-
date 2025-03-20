package storage

import (
	"database/sql"
	"db_novel_service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

func RegisterChapter(db *gorm.DB, chapter models.Chapter) (int64, error) {
	// Инициализация значений по умолчанию
	if chapter.Nodes == nil {
		chapter.Nodes = []int64{}
	}

	if chapter.Characters == nil {
		chapter.Characters = []int64{}
	}

	if chapter.UpdatedAt == nil {
		chapter.UpdatedAt = make(map[time.Time]int64)
	}

	if chapter.Name == "" {
		chapter.Name = "Новая глава"
	}

	// Маршируем JSON поля
	nodesJSON, err := json.Marshal(chapter.Nodes)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга Nodes: %w", err)
	}

	charactersJSON, err := json.Marshal(chapter.Characters)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга Characters: %w", err)
	}

	// Создаем запись с проверенными значениями
	result := db.Model(&chapter).
		Create(map[string]interface{}{
			"id":         chapter.Id,
			"name":       chapter.Name,
			"nodes":      json.RawMessage(nodesJSON),
			"characters": json.RawMessage(charactersJSON),
			"status":     chapter.Status,
			"author":     chapter.Author,
			"start_node": chapter.StartNode,
		})

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("не удалось создать запись главы")
	}

	// Обновляем временную метку
	timestamp := time.Now()
	chapter.UpdatedAt[timestamp] = chapter.Id

	// Маршируем UpdatedAt отдельно
	updatedAtJSON, err := json.Marshal(chapter.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга UpdatedAt: %w", err)
	}

	updateResult := db.Model(&chapter).
		Where("id = ?", chapter.Id).
		Updates(map[string]interface{}{
			"updated_at": json.RawMessage(updatedAtJSON),
		})

	if updateResult.Error != nil {
		return 0, fmt.Errorf("ошибка обновления временной метки: %w", updateResult.Error)
	}

	return chapter.Id, nil
}

//func SelectChapterWIthId(db *gorm.DB, id int64) (models.Chapter, error) {
//	var chapter models.Chapter
//	result := db.First(&chapter, "id = ?", id)
//	if result.RowsAffected == 0 {
//		return models.Chapter{}, errors.New("chapter data not found")
//	}
//	return chapter, nil
//}

func SelectChapterWIthId(db *gorm.DB, id int64) (models.Chapter, error) {
	var chapter models.Chapter
	query := `
        SELECT id, name, start_node,
               CAST(nodes AS TEXT) as nodes_raw,
               CAST(characters AS TEXT) as characters_raw,
               status,
               CAST(updated_at AS TEXT) as updated_at_raw,
               author
        FROM chapters
        WHERE id = $1
        LIMIT 1
    `

	row := db.Raw(query, id).Row()
	if err := row.Err(); err != nil {
		return models.Chapter{}, err
	}

	var (
		nodesRaw      string
		charactersRaw string
		updatedAtRaw  string
	)

	err := row.Scan(
		&chapter.Id,
		&chapter.Name,
		&chapter.StartNode,
		&nodesRaw,
		&charactersRaw,
		&chapter.Status,
		&updatedAtRaw,
		&chapter.Author,
	)

	log.Println(chapter)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Chapter{}, errors.New("chapter data not found")
		}
		return models.Chapter{}, err
	}

	// Десериализация JSON полей
	if err := json.Unmarshal([]byte(nodesRaw), &chapter.Nodes); err != nil {
		return models.Chapter{}, fmt.Errorf("failed to unmarshal nodes: %w", err)
	}

	if err := json.Unmarshal([]byte(charactersRaw), &chapter.Characters); err != nil {
		return models.Chapter{}, fmt.Errorf("failed to unmarshal characters: %w", err)
	}

	// Специальная обработка UpdatedAt
	var updatedAtMap map[time.Time]int64
	if err := json.Unmarshal([]byte(updatedAtRaw), &updatedAtMap); err != nil {
		return models.Chapter{}, fmt.Errorf("failed to unmarshal updated_at: %w", err)
	}
	chapter.UpdatedAt = updatedAtMap

	log.Println("chapter из бд", chapter)

	return chapter, nil
}

func UpdateChapter(db *gorm.DB, id int64, newChapter models.Chapter) (models.Chapter, error) {
	var chapter models.Chapter

	nodesJSON, err := json.Marshal(newChapter.Nodes)
	charactersJSON, err := json.Marshal(newChapter.Characters)
	updatedInfoJSON, err := json.Marshal(newChapter.UpdatedAt)

	if err != nil {
		return models.Chapter{}, err
	}

	result := db.Model(&chapter).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":       newChapter.Name,
			"StartNode":  newChapter.StartNode,
			"Nodes":      json.RawMessage(nodesJSON),
			"Characters": json.RawMessage(charactersJSON),
			"Status":     newChapter.Status,
			"UpdatedAt":  json.RawMessage(updatedInfoJSON),
		})

	if result.RowsAffected == 0 {
		return models.Chapter{}, errors.New("chapter data not update")
	}

	return chapter, nil
}

func DeleteChapter(db *gorm.DB, id int64) (int64, error) {
	var deletedChapter models.Chapter
	result := db.Where("id = ?", id).Delete(&deletedChapter)
	if result.RowsAffected == 0 {
		return 0, errors.New("payment data not update")
	}
	return result.RowsAffected, nil
}

func GetChaptersForAdmin(db *gorm.DB) ([]models.Chapter, error) {
	var chapters []models.Chapter
	query := `
        SELECT 
            id,
            name,
            start_node,
            CAST(COALESCE(nodes, '[]'::json) AS TEXT) as nodes_raw,
            CAST(COALESCE(characters, '[]'::json) AS TEXT) as characters_raw,
            status,
            CAST(updated_at AS TEXT) as updated_at_raw,
            author
        FROM chapters
    `

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch published chapters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			chapter       models.Chapter
			nodesRaw      sql.NullString
			charactersRaw sql.NullString
			updatedAtRaw  string
		)

		err = rows.Scan(
			&chapter.Id,
			&chapter.Name,
			&chapter.StartNode,
			&nodesRaw,
			&charactersRaw,
			&chapter.Status,
			&updatedAtRaw,
			&chapter.Author,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan chapter: %w", err)
		}

		// Десериализация JSON полей с обработкой NULL значений
		if nodesRaw.Valid {
			if err := json.Unmarshal([]byte(nodesRaw.String), &chapter.Nodes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal nodes: %w", err)
			}
		} else {
			chapter.Nodes = make([]int64, 0) // Инициализируйте пустым слайсом
		}

		if charactersRaw.Valid {
			if err := json.Unmarshal([]byte(charactersRaw.String), &chapter.Characters); err != nil {
				return nil, fmt.Errorf("failed to unmarshal characters: %w", err)
			}
		} else {
			chapter.Characters = make([]int64, 0) // Инициализируйте пустым слайсом
		}

		// Обработка UpdatedAt остается без изменений
		var updatedAtMap map[time.Time]int64
		if err := json.Unmarshal([]byte(updatedAtRaw), &updatedAtMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal updated_at: %w", err)
		}
		chapter.UpdatedAt = updatedAtMap

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func FindPublishedChapters(db *gorm.DB) ([]models.Chapter, error) {
	var chapters []models.Chapter
	query := `
        SELECT 
            id,
            name,
            start_node,
            CAST(nodes AS TEXT) as nodes_raw,
            CAST(characters AS TEXT) as characters_raw,
            status,
            CAST(updated_at AS TEXT) as updated_at_raw,
            author
        FROM chapters
        WHERE status = $1
    `

	rows, err := db.Raw(query, 3).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch published chapters: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			chapter       models.Chapter
			nodesRaw      string
			charactersRaw string
			updatedAtRaw  string
		)

		err = rows.Scan(
			&chapter.Id,
			&chapter.Name,
			&chapter.StartNode,
			&nodesRaw,
			&charactersRaw,
			&chapter.Status,
			&updatedAtRaw,
			&chapter.Author,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan chapter: %w", err)
		}

		// Десериализация JSON полей
		if err := json.Unmarshal([]byte(nodesRaw), &chapter.Nodes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nodes: %w", err)
		}

		if err := json.Unmarshal([]byte(charactersRaw), &chapter.Characters); err != nil {
			return nil, fmt.Errorf("failed to unmarshal characters: %w", err)
		}

		// Обработка UpdatedAt
		var updatedAtMap map[time.Time]int64
		if err := json.Unmarshal([]byte(updatedAtRaw), &updatedAtMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal updated_at: %w", err)
		}
		chapter.UpdatedAt = updatedAtMap

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
