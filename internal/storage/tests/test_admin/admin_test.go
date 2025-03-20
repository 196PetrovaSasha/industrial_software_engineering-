package test_admin

import (
	"db_novel_service/internal/models"
	"errors"
	"testing"
)

func TestRegisterAdmin(t *testing.T) {
	tests := []struct {
		name           string
		inputAdmin     models.Admin
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputAdmin: models.Admin{
				Name:  "Тестовый админ",
				Email: "test@example.com",
			},
			expectedResult: 1,
			expectedErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				RegisterAdminFunc: func(admin models.Admin) (int64, error) {
					if admin.Name == "" || admin.Email == "" {
						return 0, errors.New("admin not created")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterAdmin(tt.inputAdmin)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterAdmin() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterAdmin() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectAdminWithEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedAdmin models.Admin
		expectError   bool
	}{
		{
			name:  "Успешный поиск",
			email: "test@example.com",
			expectedAdmin: models.Admin{
				Id:               1,
				Name:             "Тестовый админ",
				Email:            "test@example.com",
				CreatedChapters:  []int64{1, 2, 3},
				RequestSent:      []int64{4, 5},
				RequestsReceived: []int64{6, 7},
			},
			expectError: false,
		},
		{
			name:          "Админ не найден",
			email:         "nonexistent@example.com",
			expectedAdmin: models.Admin{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				SelectAdminWithEmailFunc: func(email string) (models.Admin, error) {
					if email == "test@example.com" {
						return tt.expectedAdmin, nil
					}
					return models.Admin{}, errors.New("admin data not found")
				},
			}

			admin, err := mockStorage.SelectAdminWithEmail(tt.email)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectAdminWithEmail() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalAdmins(admin, tt.expectedAdmin) {
				t.Errorf("SelectAdminWithEmail() результат = %v, ожидаемый %v", admin, tt.expectedAdmin)
			}
		})
	}
}

func TestSelectAdminWithId(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		expectedAdmin models.Admin
		expectError   bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedAdmin: models.Admin{
				Id:               1,
				Name:             "Тестовый админ",
				Email:            "test@example.com",
				CreatedChapters:  []int64{1, 2, 3},
				RequestSent:      []int64{4, 5},
				RequestsReceived: []int64{6, 7},
			},
			expectError: false,
		},
		{
			name:          "Админ не найден",
			id:            999,
			expectedAdmin: models.Admin{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				SelectAdminWithIdFunc: func(id int64) (models.Admin, error) {
					if id == 1 {
						return tt.expectedAdmin, nil
					}
					return models.Admin{}, errors.New("admin data not found")
				},
			}

			admin, err := mockStorage.SelectAdminWithId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectAdminWithId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalAdmins(admin, tt.expectedAdmin) {
				t.Errorf("SelectAdminWithId() результат = %v, ожидаемый %v", admin, tt.expectedAdmin)
			}
		})
	}
}

func TestSelectAllSuperAdmins(t *testing.T) {
	tests := []struct {
		name           string
		expectedAdmins []models.Admin
		expectError    bool
	}{
		{
			name: "Найдены суперадмины",
			expectedAdmins: []models.Admin{
				{
					Id:          1,
					Name:        "Суперадмин 1",
					Email:       "super1@example.com",
					AdminStatus: 1,
				},
				{
					Id:          2,
					Name:        "Суперадмин 2",
					Email:       "super2@example.com",
					AdminStatus: 1,
				},
			},
			expectError: false,
		},
		{
			name:           "Суперадмины не найдены",
			expectedAdmins: []models.Admin{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				SelectAllSuperAdminsFunc: func() ([]models.Admin, error) {
					if len(tt.expectedAdmins) > 0 {
						return tt.expectedAdmins, nil
					}
					return nil, errors.New("no super admin found")
				},
			}

			admins, err := mockStorage.SelectAllSuperAdmins()

			if (err != nil) != tt.expectError {
				t.Errorf("SelectAllSuperAdmins() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalAdminsSlice(admins, tt.expectedAdmins) {
				t.Errorf("SelectAllSuperAdmins() результат = %v, ожидаемый %v", admins, tt.expectedAdmins)
			}
		})
	}
}

func TestUpdateAdmin(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		newAdmin      models.Admin
		expectedAdmin models.Admin
		expectError   bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newAdmin: models.Admin{
				Name:             "Обновленный админ",
				Email:            "updated@example.com",
				CreatedChapters:  []int64{10, 20, 30},
				RequestSent:      []int64{40, 50},
				RequestsReceived: []int64{60, 70},
			},
			expectedAdmin: models.Admin{
				Id:               1,
				Name:             "Обновленный админ",
				Email:            "updated@example.com",
				CreatedChapters:  []int64{10, 20, 30},
				RequestSent:      []int64{40, 50},
				RequestsReceived: []int64{60, 70},
			},
			expectError: false,
		},
		{
			name: "Админ не найден",
			id:   999,
			newAdmin: models.Admin{
				Name:  "Новый админ",
				Email: "new@example.com",
			},
			expectedAdmin: models.Admin{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockStorage{
				UpdateAdminFunc: func(id int64, newAdmin models.Admin) (models.Admin, error) {
					if id == 1 {
						return tt.expectedAdmin, nil
					}
					return models.Admin{}, errors.New("admin data not update")
				},
			}

			admin, err := mockStorage.UpdateAdmin(tt.id, tt.newAdmin)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateAdmin() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalAdmins(admin, tt.expectedAdmin) {
				t.Errorf("UpdateAdmin() результат = %v, ожидаемый %v", admin, tt.expectedAdmin)
			}
		})
	}
}

func equalAdmins(a, b models.Admin) bool {
	return a.Id == b.Id &&
		a.Name == b.Name &&
		a.Email == b.Email &&
		jsonEqual(a.CreatedChapters, b.CreatedChapters) &&
		jsonEqual(a.RequestSent, b.RequestSent) &&
		jsonEqual(a.RequestsReceived, b.RequestsReceived)
}

func jsonEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func equalAdminsSlice(a, b []models.Admin) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equalAdmins(v, b[i]) {
			return false
		}
	}
	return true
}
