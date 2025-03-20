package chapter

import (
	"db_novel_service/internal/models"
	"db_novel_service/internal/storage"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

const (
	DefaultStatus = 1
)

func CreateDefaultChapter(authorId int64, db *gorm.DB) (int64, int64, error) {

	id := generateUniqueId()

	idNode := generateUniqueId()

	newChapter := models.Chapter{
		Id:         id,
		Author:     authorId,
		Status:     DefaultStatus,
		UpdatedAt:  map[time.Time]int64{time.Now(): id},
		Characters: []int64{},
	}

	_, err := storage.RegisterChapter(db, newChapter)

	log.Println(err)

	newNode := models.Node{
		Id:   idNode,
		Slug: "",
	}

	nodeId, err := CreateNode(id, newNode.Slug, db)

	log.Println("id начального узла", nodeId)

	chapter, err := storage.SelectChapterWIthId(db, id)

	chapter.StartNode = nodeId

	_, err = storage.UpdateChapter(db, id, chapter)

	if err != nil {
		log.Println("ошибка инициализации начального узла")
		return 0, 0, err
	}

	return id, nodeId, nil
}

func generateUniqueId() int64 {
	// Получаем текущее время в миллисекундах (48 бит)
	timestamp := time.Now().UnixMilli()

	// Генерируем 16 случайных бит
	random := rand.Int31n(1 << 16)

	// Объединяем timestamp и random в 64-битное число
	return (int64(timestamp) << 16) | int64(random)
}

func CreateNode(chapterId int64, slug string, db *gorm.DB) (int64, error) {

	id := generateUniqueId()

	newNode := models.Node{
		Id:        id,
		Slug:      slug,
		ChapterId: chapterId,
		Events:    map[int]models.Event{},
		Branching: models.Branching{},
		End:       models.EndInfo{},
		Comment:   " ",
	}

	_, err := storage.RegisterNode(db, newNode)

	if err != nil {
		return 0, err
	}

	chapter, err := storage.SelectChapterWIthId(db, chapterId)

	if err != nil {
		return 0, err
	}

	chapter.Nodes = append(chapter.Nodes, id)

	_, err = storage.UpdateChapter(db, chapterId, chapter)

	if err != nil {
		return 0, err
	}

	return id, nil
}
