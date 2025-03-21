package test_node

import (
	"errors"
	"testing"
	"vn/internal/models"
)

func TestRegisterNode(t *testing.T) {
	tests := []struct {
		name           string
		inputNode      models.Node
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputNode: models.Node{
				Slug: "test-slug",
				Events: map[int]models.Event{
					1: {},
				},
				Branching:  models.Branching{},
				End:        models.EndInfo{},
				Music:      1,
				Background: 2,
				Comment:    "Test comment",
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Ошибка валидации",
			inputNode: models.Node{
				Slug: "", // пустой slug
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockNodeStorage{
				RegisterNodeFunc: func(node models.Node) (int64, error) {
					if node.Slug == "" {
						return 0, errors.New("не удалось создать запись узла")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterNode(tt.inputNode)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterNode() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterNode() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectNodeWithId(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		expectedNode *models.Node
		expectError  bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedNode: &models.Node{
				Id:   1,
				Slug: "test-slug",
				Events: map[int]models.Event{
					1: {},
				},
				Branching:  models.Branching{},
				End:        models.EndInfo{},
				Music:      1,
				Background: 2,
				Comment:    "Test comment",
			},
			expectError: false,
		},
		{
			name:         "Узел не найден",
			id:           999,
			expectedNode: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockNodeStorage{
				SelectNodeWithIdFunc: func(id int64) (*models.Node, error) {
					if id == 1 {
						return tt.expectedNode, nil
					}
					return nil, errors.New("node data not found")
				},
			}

			node, err := mockStorage.SelectNodeWithId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectNodeWithId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalNodes(node, tt.expectedNode) {
				t.Errorf("SelectNodeWithId() результат = %v, ожидаемый %v", node, tt.expectedNode)
			}
		})
	}
}

func TestUpdateNode(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		newNode      models.Node
		expectedNode models.Node
		expectError  bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newNode: models.Node{
				Slug: "updated-slug",
				Events: map[int]models.Event{
					1: {},
				},
				Branching:  models.Branching{},
				End:        models.EndInfo{},
				Music:      2,
				Background: 3,
				Comment:    "Updated comment",
			},
			expectedNode: models.Node{
				Id:   1,
				Slug: "updated-slug",
				Events: map[int]models.Event{
					1: {},
				},
				Branching:  models.Branching{},
				End:        models.EndInfo{},
				Music:      2,
				Background: 3,
				Comment:    "Updated comment",
			},
			expectError: false,
		},
		{
			name: "Узел не найден",
			id:   999,
			newNode: models.Node{
				Slug: "new-slug",
			},
			expectedNode: models.Node{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockNodeStorage{
				UpdateNodeFunc: func(id int64, newNode models.Node) (models.Node, error) {
					if id == 1 {
						return tt.expectedNode, nil
					}
					return models.Node{}, errors.New("node data not update")
				},
			}

			node, err := mockStorage.UpdateNode(tt.id, tt.newNode)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateNode() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalNodes(&node, &tt.expectedNode) {
				t.Errorf("UpdateNode() результат = %v, ожидаемый %v", node, tt.expectedNode)
			}
		})
	}
}

func TestDeleteNode(t *testing.T) {
	tests := []struct {
		name           string
		id             int64
		expectedResult int64
		expectError    bool
	}{
		{
			name:           "Успешное удаление",
			id:             1,
			expectedResult: 1,
			expectError:    false,
		},
		{
			name:           "Узел не найден",
			id:             999,
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockNodeStorage{
				DeleteNodeFunc: func(id int64) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("запись не найдена")
				},
			}

			result, err := mockStorage.DeleteNode(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("DeleteNode() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeleteNode() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGetNodeById(t *testing.T) {
	tests := []struct {
		name         string
		id           int64
		expectedNode models.Node
		expectError  bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedNode: models.Node{
				Id:   1,
				Slug: "test-slug",
				Events: map[int]models.Event{
					1: {},
				},
				Branching:  models.Branching{},
				End:        models.EndInfo{},
				Music:      1,
				Background: 2,
				Comment:    "Test comment",
			},
			expectError: false,
		},
		{
			name:         "Узел не найден",
			id:           999,
			expectedNode: models.Node{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockNodeStorage{
				GetNodeByIdFunc: func(id int64) (models.Node, error) {
					if id == 1 {
						return tt.expectedNode, nil
					}
					return models.Node{}, errors.New("failed to fetch node")
				},
			}

			node, err := mockStorage.GetNodeById(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("GetNodeById() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalNodes(&node, &tt.expectedNode) {
				t.Errorf("GetNodeById() результат = %v, ожидаемый %v", node, tt.expectedNode)
			}
		})
	}
}

func equalNodes(a, b *models.Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Id == b.Id &&
		a.Slug == b.Slug &&
		a.Music == b.Music &&
		a.Background == b.Background &&
		a.Comment == b.Comment
}
