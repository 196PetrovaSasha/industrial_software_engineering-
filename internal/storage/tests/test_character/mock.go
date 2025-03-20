package test_character

import (
	"db_novel_service/internal/models"
	"errors"
)

// CharacterStorage реализует интерфейс для работы с персонажами
type CharacterStorage interface {
	RegisterCharacter(character models.Character) (int64, error)
	SelectCharacterWithId(id int64) (models.Character, error)
	SelectCharacters() ([]models.Character, error)
	UpdateCharacter(id int64, newCharacter models.Character) (models.Character, error)
	DeleteCharacter(id int64) (int64, error)
}

// MockCharacterStorage реализует интерфейс CharacterStorage для тестирования
type MockCharacterStorage struct {
	RegisterCharacterFunc     func(character models.Character) (int64, error)
	SelectCharacterWithIdFunc func(id int64) (models.Character, error)
	SelectCharactersFunc      func() ([]models.Character, error)
	UpdateCharacterFunc       func(id int64, newCharacter models.Character) (models.Character, error)
	DeleteCharacterFunc       func(id int64) (int64, error)
}

// RegisterCharacter реализация для MockCharacterStorage
func (m *MockCharacterStorage) RegisterCharacter(character models.Character) (int64, error) {
	if m.RegisterCharacterFunc != nil {
		return m.RegisterCharacterFunc(character)
	}
	return 0, errors.New("RegisterCharacter не реализован")
}

// SelectCharacterWithId реализация для MockCharacterStorage
func (m *MockCharacterStorage) SelectCharacterWithId(id int64) (models.Character, error) {
	if m.SelectCharacterWithIdFunc != nil {
		return m.SelectCharacterWithIdFunc(id)
	}
	return models.Character{}, errors.New("SelectCharacterWithId не реализован")
}

// SelectCharacters реализация для MockCharacterStorage
func (m *MockCharacterStorage) SelectCharacters() ([]models.Character, error) {
	if m.SelectCharactersFunc != nil {
		return m.SelectCharactersFunc()
	}
	return nil, errors.New("SelectCharacters не реализован")
}

// UpdateCharacter реализация для MockCharacterStorage
func (m *MockCharacterStorage) UpdateCharacter(id int64, newCharacter models.Character) (models.Character, error) {
	if m.UpdateCharacterFunc != nil {
		return m.UpdateCharacterFunc(id, newCharacter)
	}
	return models.Character{}, errors.New("UpdateCharacter не реализован")
}

// DeleteCharacter реализация для MockCharacterStorage
func (m *MockCharacterStorage) DeleteCharacter(id int64) (int64, error) {
	if m.DeleteCharacterFunc != nil {
		return m.DeleteCharacterFunc(id)
	}
	return 0, errors.New("DeleteCharacter не реализован")
}
