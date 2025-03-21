package test_media

import (
	"errors"
	"vn/internal/models"
)

// MediaStorage реализует интерфейс для работы с медиа
type MediaStorage interface {
	RegisterMedia(media models.Media) (int64, error)
	SelectMediaWIthId(id int64) (models.Media, error)
	SelectMedia() ([]models.Media, error)
	DeleteMedia(id int64) (int64, error)
}

// MockMediaStorage реализует интерфейс MediaStorage для тестирования
type MockMediaStorage struct {
	RegisterMediaFunc     func(media models.Media) (int64, error)
	SelectMediaWIthIdFunc func(id int64) (models.Media, error)
	SelectMediaFunc       func() ([]models.Media, error)
	DeleteMediaFunc       func(id int64) (int64, error)
}

// RegisterMedia реализация для MockMediaStorage
func (m *MockMediaStorage) RegisterMedia(media models.Media) (int64, error) {
	if m.RegisterMediaFunc != nil {
		return m.RegisterMediaFunc(media)
	}
	return 0, errors.New("RegisterMedia не реализован")
}

// SelectMediaWIthId реализация для MockMediaStorage
func (m *MockMediaStorage) SelectMediaWIthId(id int64) (models.Media, error) {
	if m.SelectMediaWIthIdFunc != nil {
		return m.SelectMediaWIthIdFunc(id)
	}
	return models.Media{}, errors.New("SelectMediaWIthId не реализован")
}

// SelectMedia реализация для MockMediaStorage
func (m *MockMediaStorage) SelectMedia() ([]models.Media, error) {
	if m.SelectMediaFunc != nil {
		return m.SelectMediaFunc()
	}
	return nil, errors.New("SelectMedia не реализован")
}

// DeleteMedia реализация для MockMediaStorage
func (m *MockMediaStorage) DeleteMedia(id int64) (int64, error) {
	if m.DeleteMediaFunc != nil {
		return m.DeleteMediaFunc(id)
	}
	return 0, errors.New("DeleteMedia не реализован")
}
