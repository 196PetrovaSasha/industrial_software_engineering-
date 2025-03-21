package test_player

import (
	"errors"
	"vn/internal/models"
)

// PlayerStorage реализует интерфейс для работы с игроками
type PlayerStorage interface {
	RegisterPlayer(player models.Player) (int64, error)
	SelectPlayerWIthEmail(email string) (models.Player, error)
	SelectPlayerWIthId(id int64) (models.Player, error)
	UpdatePlayer(id int64, newPlayer models.Player) (models.Player, error)
	DeletePlayer(email string) (int64, error)
}

// MockPlayerStorage реализует интерфейс PlayerStorage для тестирования
type MockPlayerStorage struct {
	RegisterPlayerFunc        func(player models.Player) (int64, error)
	SelectPlayerWIthEmailFunc func(email string) (models.Player, error)
	SelectPlayerWIthIdFunc    func(id int64) (models.Player, error)
	UpdatePlayerFunc          func(id int64, newPlayer models.Player) (models.Player, error)
	DeletePlayerFunc          func(email string) (int64, error)
}

func (m *MockPlayerStorage) RegisterPlayer(player models.Player) (int64, error) {
	if m.RegisterPlayerFunc != nil {
		return m.RegisterPlayerFunc(player)
	}
	return 0, errors.New("RegisterPlayer не реализован")
}

func (m *MockPlayerStorage) SelectPlayerWIthEmail(email string) (models.Player, error) {
	if m.SelectPlayerWIthEmailFunc != nil {
		return m.SelectPlayerWIthEmailFunc(email)
	}
	return models.Player{}, errors.New("SelectPlayerWIthEmail не реализован")
}

func (m *MockPlayerStorage) SelectPlayerWIthId(id int64) (models.Player, error) {
	if m.SelectPlayerWIthIdFunc != nil {
		return m.SelectPlayerWIthIdFunc(id)
	}
	return models.Player{}, errors.New("SelectPlayerWIthId не реализован")
}

func (m *MockPlayerStorage) UpdatePlayer(id int64, newPlayer models.Player) (models.Player, error) {
	if m.UpdatePlayerFunc != nil {
		return m.UpdatePlayerFunc(id, newPlayer)
	}
	return models.Player{}, errors.New("UpdatePlayer не реализован")
}

func (m *MockPlayerStorage) DeletePlayer(email string) (int64, error) {
	if m.DeletePlayerFunc != nil {
		return m.DeletePlayerFunc(email)
	}
	return 0, errors.New("DeletePlayer не реализован")
}
