package test_request

import (
	"db_novel_service/internal/models"
	"errors"
	"fmt"
	"testing"
)

func TestRegisterRequest(t *testing.T) {
	tests := []struct {
		name           string
		inputRequest   models.Request
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputRequest: models.Request{
				Type:               1,
				Status:             0,
				RequestingAdmin:    1,
				RequestedChapterId: 1,
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Ошибка регистрации",
			inputRequest: models.Request{
				Type: -1, // неверный тип
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockRequestStorage{
				RegisterRequestFunc: func(request models.Request) (int64, error) {
					if request.Type < 0 || request.Type > 3 {
						return 0, errors.New("неверный тип запроса")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterRequest(tt.inputRequest)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterRequest() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterRequest() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectRequestWithId(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		expectedRequest *models.Request
		expectError     bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedRequest: &models.Request{
				Id:                 1,
				Type:               1,
				Status:             0,
				RequestingAdmin:    1,
				RequestedChapterId: 1,
			},
			expectError: false,
		},
		{
			name:            "Запрос не найден",
			id:              999,
			expectedRequest: nil,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockRequestStorage{
				SelectRequestWithIdFunc: func(id int64) (*models.Request, error) {
					if id == 1 {
						return tt.expectedRequest, nil
					}
					return nil, fmt.Errorf("запрос с ID %d не найден", id)
				},
			}

			request, err := mockStorage.SelectRequestWithId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectRequestWithId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalRequests(request, tt.expectedRequest) {
				t.Errorf("SelectRequestWithId() результат = %v, ожидаемый %v", request, tt.expectedRequest)
			}
		})
	}
}

func TestUpdateRequest(t *testing.T) {
	tests := []struct {
		name           string
		id             int64
		newRequest     models.Request
		expectedResult int64
		expectError    bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newRequest: models.Request{
				Type:               2,
				Status:             1,
				RequestingAdmin:    1,
				RequestedChapterId: 2,
			},
			expectedResult: 1,
			expectError:    false,
		},
		{
			name: "Запрос не найден",
			id:   999,
			newRequest: models.Request{
				Type: 1,
			},
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockRequestStorage{
				UpdateRequestFunc: func(id int64, newRequest models.Request) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("запрос не обновлен")
				},
			}

			result, err := mockStorage.UpdateRequest(tt.id, tt.newRequest)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateRequest() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("UpdateRequest() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestDeleteRequest(t *testing.T) {
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
			name:           "Запрос не найден",
			id:             999,
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockRequestStorage{
				DeleteRequestFunc: func(id int64) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("запрос не удален")
				},
			}

			result, err := mockStorage.DeleteRequest(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("DeleteRequest() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeleteRequest() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGetAllRequests(t *testing.T) {
	tests := []struct {
		name             string
		expectedRequests []models.Request
		expectError      bool
	}{
		{
			name: "Успешное получение",
			expectedRequests: []models.Request{
				{
					Id:                 1,
					Type:               1,
					Status:             0,
					RequestingAdmin:    1,
					RequestedChapterId: 1,
				},
				{
					Id:                 2,
					Type:               2,
					Status:             1,
					RequestingAdmin:    2,
					RequestedChapterId: 2,
				},
			},
			expectError: false,
		},
		{
			name:             "Ошибка получения",
			expectedRequests: []models.Request{},
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockRequestStorage{
				GetAllRequestsFunc: func() ([]models.Request, error) {
					if len(tt.expectedRequests) > 0 {
						return tt.expectedRequests, nil
					}
					return nil, errors.New("запросы не найдены")
				},
			}

			requests, err := mockStorage.GetAllRequests()

			if (err != nil) != tt.expectError {
				t.Errorf("GetAllRequests() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalRequestsSlice(requests, tt.expectedRequests) {
				t.Errorf("GetAllRequests() результат = %v, ожидаемый %v", requests, tt.expectedRequests)
			}
		})
	}
}

func equalRequests(a, b *models.Request) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Id == b.Id &&
		a.Type == b.Type &&
		a.Status == b.Status &&
		a.RequestingAdmin == b.RequestingAdmin &&
		a.RequestedChapterId == b.RequestedChapterId
}

func equalRequestsSlice(a, b []models.Request) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equalRequests(&v, &b[i]) {
			return false
		}
	}
	return true
}
