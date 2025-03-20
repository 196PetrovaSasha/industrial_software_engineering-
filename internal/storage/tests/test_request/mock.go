package test_request

import (
	"db_novel_service/internal/models"
	"errors"
)

// RequestStorage реализует интерфейс для работы с запросами
type RequestStorage interface {
	RegisterRequest(request models.Request) (int64, error)
	SelectRequestWithId(id int64) (*models.Request, error)
	UpdateRequest(id int64, newRequest models.Request) (int64, error)
	DeleteRequest(id int64) (int64, error)
	GetAllRequests() ([]models.Request, error)
}

// MockRequestStorage реализует интерфейс RequestStorage для тестирования
type MockRequestStorage struct {
	RegisterRequestFunc     func(request models.Request) (int64, error)
	SelectRequestWithIdFunc func(id int64) (*models.Request, error)
	UpdateRequestFunc       func(id int64, newRequest models.Request) (int64, error)
	DeleteRequestFunc       func(id int64) (int64, error)
	GetAllRequestsFunc      func() ([]models.Request, error)
}

func (m *MockRequestStorage) RegisterRequest(request models.Request) (int64, error) {
	if m.RegisterRequestFunc != nil {
		return m.RegisterRequestFunc(request)
	}
	return 0, errors.New("RegisterRequest не реализован")
}

func (m *MockRequestStorage) SelectRequestWithId(id int64) (*models.Request, error) {
	if m.SelectRequestWithIdFunc != nil {
		return m.SelectRequestWithIdFunc(id)
	}
	return nil, errors.New("SelectRequestWithId не реализован")
}

func (m *MockRequestStorage) UpdateRequest(id int64, newRequest models.Request) (int64, error) {
	if m.UpdateRequestFunc != nil {
		return m.UpdateRequestFunc(id, newRequest)
	}
	return 0, errors.New("UpdateRequest не реализован")
}

func (m *MockRequestStorage) DeleteRequest(id int64) (int64, error) {
	if m.DeleteRequestFunc != nil {
		return m.DeleteRequestFunc(id)
	}
	return 0, errors.New("DeleteRequest не реализован")
}

func (m *MockRequestStorage) GetAllRequests() ([]models.Request, error) {
	if m.GetAllRequestsFunc != nil {
		return m.GetAllRequestsFunc()
	}
	return nil, errors.New("GetAllRequests не реализован")
}
