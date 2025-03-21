package test_chapter

import (
	"errors"
	"vn/internal/models"
)

// ChapterStorage реализует интерфейс для работы с главами
type ChapterStorage interface {
	RegisterChapter(chapter models.Chapter) (int64, error)
	SelectChapterWithId(id int64) (models.Chapter, error)
	UpdateChapter(id int64, newChapter models.Chapter) (models.Chapter, error)
	DeleteChapter(id int64) (int64, error)
	GetChaptersForAdmin() ([]models.Chapter, error)
	FindPublishedChapters() ([]models.Chapter, error)
}

// MockChapterStorage реализует интерфейс ChapterStorage для тестирования
type MockChapterStorage struct {
	RegisterChapterFunc       func(chapter models.Chapter) (int64, error)
	SelectChapterWithIdFunc   func(id int64) (models.Chapter, error)
	UpdateChapterFunc         func(id int64, newChapter models.Chapter) (models.Chapter, error)
	DeleteChapterFunc         func(id int64) (int64, error)
	GetChaptersForAdminFunc   func() ([]models.Chapter, error)
	FindPublishedChaptersFunc func() ([]models.Chapter, error)
}

// RegisterChapter реализация для MockChapterStorage
func (m *MockChapterStorage) RegisterChapter(chapter models.Chapter) (int64, error) {
	if m.RegisterChapterFunc != nil {
		return m.RegisterChapterFunc(chapter)
	}
	return 0, errors.New("RegisterChapter не реализован")
}

// SelectChapterWithId реализация для MockChapterStorage
func (m *MockChapterStorage) SelectChapterWithId(id int64) (models.Chapter, error) {
	if m.SelectChapterWithIdFunc != nil {
		return m.SelectChapterWithIdFunc(id)
	}
	return models.Chapter{}, errors.New("SelectChapterWithId не реализован")
}

// UpdateChapter реализация для MockChapterStorage
func (m *MockChapterStorage) UpdateChapter(id int64, newChapter models.Chapter) (models.Chapter, error) {
	if m.UpdateChapterFunc != nil {
		return m.UpdateChapterFunc(id, newChapter)
	}
	return models.Chapter{}, errors.New("UpdateChapter не реализован")
}

// DeleteChapter реализация для MockChapterStorage
func (m *MockChapterStorage) DeleteChapter(id int64) (int64, error) {
	if m.DeleteChapterFunc != nil {
		return m.DeleteChapterFunc(id)
	}
	return 0, errors.New("DeleteChapter не реализован")
}

// GetChaptersForAdmin реализация для MockChapterStorage
func (m *MockChapterStorage) GetChaptersForAdmin() ([]models.Chapter, error) {
	if m.GetChaptersForAdminFunc != nil {
		return m.GetChaptersForAdminFunc()
	}
	return nil, errors.New("GetChaptersForAdmin не реализован")
}

// FindPublishedChapters реализация для MockChapterStorage
func (m *MockChapterStorage) FindPublishedChapters() ([]models.Chapter, error) {
	if m.FindPublishedChaptersFunc != nil {
		return m.FindPublishedChaptersFunc()
	}
	return nil, errors.New("FindPublishedChapters не реализован")
}
