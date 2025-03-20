package storage

import (
	"db_novel_service/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func RegisterNode(db *gorm.DB, node models.Node) (int64, error) {
	// Инициализация значений по умолчанию...
	if node.Events == nil {
		node.Events = map[int]models.Event{}
	}

	node.Branching = models.Branching{}

	node.End = models.EndInfo{}

	// Проверяем и устанавливаем значения по умолчанию для простых полей
	if node.Slug == "" {
		node.Slug = "default-slug"
	}

	if node.Music == 0 {
		node.Music = 1
	}

	if node.Background == 0 {
		node.Background = 1
	}

	// Маршируем JSON поля
	eventsJSON, err := json.Marshal(node.Events)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга Events: %w", err)
	}

	log.Printf("Events JSON: %s", eventsJSON)

	branchingJSON, err := json.Marshal(node.Branching)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга Branching: %w", err)
	}

	endInfoJSON, err := json.Marshal(node.End)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга End: %w", err)
	}

	// Создаем запись с проверенными значениями
	result := db.Model(&node).
		Create(map[string]interface{}{
			"id":         node.Id,
			"slug":       node.Slug,
			"events":     json.RawMessage(eventsJSON),
			"chapter_id": node.ChapterId,
			"music":      node.Music,
			"background": node.Background,
			"branching":  json.RawMessage(branchingJSON),
			"end_info":   json.RawMessage(endInfoJSON),
			"comment":    node.Comment,
		})

	if result.RowsAffected == 0 {
		return 0, fmt.Errorf("не удалось создать запись узла")
	}

	return node.Id, nil
}

//func RegisterNode(db *gorm.DB, node models.Node) (int64, error) {
//	// Создаем новую структуру с простыми полями
//	newNode := models.Node{
//		Id:         node.Id,
//		Slug:       node.Slug,
//		ChapterId:  node.ChapterId,
//		Music:      node.Music,
//		Background: node.Background,
//	}
//
//	// Копируем простые поля
//	newNode.Slug = node.Slug
//	newNode.ChapterId = node.ChapterId
//	newNode.Music = node.Music
//	newNode.Background = node.Background
//
//	// Копируем Events напрямую
//	newNode.Events = node.Events
//
//	// Копируем Branching напрямую
//	newNode.Branching = node.Branching
//
//	// Копируем EndInfo напрямую
//	newNode.End = node.End
//
//	result := db.Create(&newNode)
//
//	if result.RowsAffected == 0 {
//		return 0, errors.New("node not created")
//	}
//
//	return result.RowsAffected, nil
//}

//func SelectNodeWIthId(db *gorm.DB, id int64) (models.Node, error) {
//	var node models.Node
//	result := db.First(&node, "id = ?", id)
//	if result.RowsAffected == 0 {
//		return models.Node{}, errors.New("node data not found")
//	}
//	return node, nil
//}

func SelectNodeWIthId(db *gorm.DB, nodeId int64) (*models.Node, error) {
	var nodes []models.Node

	// Используем raw SQL с явной обработкой JSON
	query := `
        SELECT 
            id,
            slug,
            chapter_id,
            music,
            background,
            CAST(events AS TEXT) as events_raw,
            CAST(branching AS TEXT) as branching_raw,
            CAST(end_info AS TEXT) as end_raw,
            comment
        FROM nodes
    `

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var n models.Node
		var (
			eventsRaw    string
			branchingRaw string
			endRaw       string
		)

		// Читаем все поля, включая сырые данные JSON
		err = rows.Scan(
			&n.Id,
			&n.Slug,
			&n.ChapterId,
			&n.Music,
			&n.Background,
			&eventsRaw,
			&branchingRaw,
			&endRaw,
			&n.Comment,
		)
		if err != nil {
			return nil, err
		}

		// Десериализуем JSON в структуры
		if err := json.Unmarshal([]byte(eventsRaw), &n.Events); err != nil {
			return nil, fmt.Errorf("failed to unmarshal events: %w", err)
		}

		if err := json.Unmarshal([]byte(branchingRaw), &n.Branching); err != nil {
			return nil, fmt.Errorf("failed to unmarshal branching: %w", err)
		}

		if err := json.Unmarshal([]byte(endRaw), &n.End); err != nil {
			return nil, fmt.Errorf("failed to unmarshal end info: %w", err)
		}

		nodes = append(nodes, n)
	}

	for _, i := range nodes {
		if i.Id == nodeId {
			return &i, err
		}
	}

	return nil, rows.Err()
}

//func UpdateNode(db *gorm.DB, id int64, newNode models.Node) (models.Node, error) {
//	var node models.Node
//	result := db.Model(&node).Where("id = ?", id).Updates(newNode)
//	if result.RowsAffected == 0 {
//		return models.Node{}, errors.New("node data not update")
//	}
//	return node, nil
//}

func UpdateNode(db *gorm.DB, id int64, newNode models.Node) (models.Node, error) {
	var node models.Node
	// Сериализуем JSON поля
	eventsJSON, err := json.Marshal(newNode.Events)
	if err != nil {
		return models.Node{}, err
	}

	log.Printf("%s, events", eventsJSON)

	branchingJSON, err := json.Marshal(newNode.Branching)
	if err != nil {
		return models.Node{}, err
	}

	endInfoJSON, err := json.Marshal(newNode.End)
	if err != nil {
		return models.Node{}, err
	}

	// Обновляем данные с использованием raw SQL
	result := db.Model(&node).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"slug":       newNode.Slug,
			"events":     json.RawMessage(eventsJSON),
			"chapter_id": newNode.ChapterId,
			"music":      newNode.Music,
			"background": newNode.Background,
			"branching":  json.RawMessage(branchingJSON),
			"end_info":   json.RawMessage(endInfoJSON),
		})

	if result.RowsAffected == 0 {
		return models.Node{}, errors.New("node data not update")
	}

	return node, nil
}

func DeleteNode(db *gorm.DB, id int64) (int64, error) {
	var deletedNode models.Node
	result := db.Where("id = ?", id).Delete(&deletedNode)
	if result.Error != nil {
		return 0, fmt.Errorf("ошибка при удалении узла: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.New("запись не найдена")
	}
	return result.RowsAffected, nil
}

func GetNodeById(db *gorm.DB, id int64) (models.Node, error) {
	var node models.Node
	result := db.Where("id = ?", id).Find(&node)

	if result.Error != nil {
		return models.Node{}, fmt.Errorf("failed to fetch node: %v", result.Error)
	}

	return node, nil
}
