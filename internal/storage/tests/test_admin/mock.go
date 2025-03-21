package test_admin

import (
	"errors"
	"vn/internal/models"
)

// MockStorage реализует интерфейс Storage для тестирования
type MockStorage struct {
	RegisterAdminFunc        func(admin models.Admin) (int64, error)
	SelectAdminWithEmailFunc func(email string) (models.Admin, error)
	SelectAdminWithIdFunc    func(id int64) (models.Admin, error)
	SelectAllSuperAdminsFunc func() ([]models.Admin, error)
	UpdateAdminFunc          func(id int64, newAdmin models.Admin) (models.Admin, error)
}

// RegisterAdmin реализация для MockStorage
func (m *MockStorage) RegisterAdmin(admin models.Admin) (int64, error) {
	if m.RegisterAdminFunc != nil {
		return m.RegisterAdminFunc(admin)
	}
	return 0, errors.New("RegisterAdmin не реализован")
}

// SelectAdminWithEmail реализация для MockStorage
func (m *MockStorage) SelectAdminWithEmail(email string) (models.Admin, error) {
	if m.SelectAdminWithEmailFunc != nil {
		return m.SelectAdminWithEmailFunc(email)
	}
	return models.Admin{}, errors.New("SelectAdminWithEmail не реализован")
}

// SelectAdminWithId реализация для MockStorage
func (m *MockStorage) SelectAdminWithId(id int64) (models.Admin, error) {
	if m.SelectAdminWithIdFunc != nil {
		return m.SelectAdminWithIdFunc(id)
	}
	return models.Admin{}, errors.New("SelectAdminWithId не реализован")
}

// SelectAllSuperAdmins реализация для MockStorage
func (m *MockStorage) SelectAllSuperAdmins() ([]models.Admin, error) {
	if m.SelectAllSuperAdminsFunc != nil {
		return m.SelectAllSuperAdminsFunc()
	}
	return nil, errors.New("SelectAllSuperAdmins не реализован")
}

// UpdateAdmin реализация для MockStorage
func (m *MockStorage) UpdateAdmin(id int64, newAdmin models.Admin) (models.Admin, error) {
	if m.UpdateAdminFunc != nil {
		return m.UpdateAdminFunc(id, newAdmin)
	}
	return models.Admin{}, errors.New("UpdateAdmin не реализован")
}
