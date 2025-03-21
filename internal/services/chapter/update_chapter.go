package chapter

import (
	"gorm.io/gorm"
	"time"
	"vn/internal/storage"
)

func UpdateChapter(
	id int64,
	name string,
	nodes []int64,
	characters []int64,
	updateAuthorId int64,
	startNode int64,
	status int,
	db *gorm.DB,
) error {
	chapter, err := storage.SelectChapterWIthId(db, id)

	if err != nil {
		return err
	}

	newChapter := chapter

	if newChapter.UpdatedAt == nil {
		newChapter.UpdatedAt = make(map[time.Time]int64)
	}

	newChapter.UpdatedAt[time.Now()] = updateAuthorId

	if name != "" {
		newChapter.Name = name
	}

	if nodes != nil {
		newChapter.Nodes = nodes
	}

	if characters != nil {
		newChapter.Characters = characters
	}

	if startNode != 0 {
		newChapter.StartNode = startNode
	}

	if status != 0 {
		newChapter.Status = status
	}

	_, err = storage.UpdateChapter(db, id, newChapter)

	return err
}
