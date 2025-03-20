package test_node

import (
	"db_novel_service/internal/models"
	"errors"
)

// NodeStorage реализует интерфейс для работы с узлами
type NodeStorage interface {
	RegisterNode(node models.Node) (int64, error)
	SelectNodeWithId(id int64) (*models.Node, error)
	UpdateNode(id int64, newNode models.Node) (models.Node, error)
	DeleteNode(id int64) (int64, error)
	GetNodeById(id int64) (models.Node, error)
}

// MockNodeStorage реализует интерфейс NodeStorage для тестирования
type MockNodeStorage struct {
	RegisterNodeFunc     func(node models.Node) (int64, error)
	SelectNodeWithIdFunc func(id int64) (*models.Node, error)
	UpdateNodeFunc       func(id int64, newNode models.Node) (models.Node, error)
	DeleteNodeFunc       func(id int64) (int64, error)
	GetNodeByIdFunc      func(id int64) (models.Node, error)
}

// RegisterNode реализация для MockNodeStorage
func (m *MockNodeStorage) RegisterNode(node models.Node) (int64, error) {
	if m.RegisterNodeFunc != nil {
		return m.RegisterNodeFunc(node)
	}
	return 0, errors.New("RegisterNode не реализован")
}

// SelectNodeWithId реализация для MockNodeStorage
func (m *MockNodeStorage) SelectNodeWithId(id int64) (*models.Node, error) {
	if m.SelectNodeWithIdFunc != nil {
		return m.SelectNodeWithIdFunc(id)
	}
	return nil, errors.New("SelectNodeWithId не реализован")
}

// UpdateNode реализация для MockNodeStorage
func (m *MockNodeStorage) UpdateNode(id int64, newNode models.Node) (models.Node, error) {
	if m.UpdateNodeFunc != nil {
		return m.UpdateNodeFunc(id, newNode)
	}
	return models.Node{}, errors.New("UpdateNode не реализован")
}

// DeleteNode реализация для MockNodeStorage
func (m *MockNodeStorage) DeleteNode(id int64) (int64, error) {
	if m.DeleteNodeFunc != nil {
		return m.DeleteNodeFunc(id)
	}
	return 0, errors.New("DeleteNode не реализован")
}

// GetNodeById реализация для MockNodeStorage
func (m *MockNodeStorage) GetNodeById(id int64) (models.Node, error) {
	if m.GetNodeByIdFunc != nil {
		return m.GetNodeByIdFunc(id)
	}
	return models.Node{}, errors.New("GetNodeById не реализован")
}
