package test_character

import (
	"db_novel_service/internal/models"
	"errors"
	"testing"
)

func TestRegisterCharacter(t *testing.T) {
	tests := []struct {
		name           string
		inputCharacter models.Character
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputCharacter: models.Character{
				Name:  "Тестовый персонаж",
				Slug:  "test-character",
				Color: "blue",
				Emotions: map[int64]int64{
					1: 1,
					2: 2,
				},
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Ошибка регистрации",
			inputCharacter: models.Character{
				Name: "", // пустое имя
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockCharacterStorage{
				RegisterCharacterFunc: func(character models.Character) (int64, error) {
					if character.Name == "" {
						return 0, errors.New("character not created")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterCharacter(tt.inputCharacter)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterCharacter() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterCharacter() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectCharacterWithId(t *testing.T) {
	tests := []struct {
		name              string
		id                int64
		expectedCharacter models.Character
		expectError       bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedCharacter: models.Character{
				Id:    1,
				Name:  "Тестовый персонаж",
				Slug:  "test-character",
				Color: "blue",
				Emotions: map[int64]int64{
					1: 1,
					2: 2,
				},
			},
			expectError: false,
		},
		{
			name:              "Персонаж не найден",
			id:                999,
			expectedCharacter: models.Character{},
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockCharacterStorage{
				SelectCharacterWithIdFunc: func(id int64) (models.Character, error) {
					if id == 1 {
						return tt.expectedCharacter, nil
					}
					return models.Character{}, errors.New("character not found")
				},
			}

			character, err := mockStorage.SelectCharacterWithId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectCharacterWithId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalCharacters(character, tt.expectedCharacter) {
				t.Errorf("SelectCharacterWithId() результат = %v, ожидаемый %v", character, tt.expectedCharacter)
			}
		})
	}
}

func TestSelectCharacters(t *testing.T) {
	tests := []struct {
		name               string
		expectedCharacters []models.Character
		expectError        bool
	}{
		{
			name: "Успешный поиск",
			expectedCharacters: []models.Character{
				{
					Id:    1,
					Name:  "Персонаж 1",
					Slug:  "character-1",
					Color: "red",
					Emotions: map[int64]int64{
						1: 1,
						2: 2,
					},
				},
				{
					Id:    2,
					Name:  "Персонаж 2",
					Slug:  "character-2",
					Color: "green",
					Emotions: map[int64]int64{
						1: 1,
						2: 2,
					},
				},
			},
			expectError: false,
		},
		{
			name:               "Ошибка получения персонажей",
			expectedCharacters: []models.Character{},
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockCharacterStorage{
				SelectCharactersFunc: func() ([]models.Character, error) {
					if len(tt.expectedCharacters) > 0 {
						return tt.expectedCharacters, nil
					}
					return nil, errors.New("characters not found")
				},
			}

			characters, err := mockStorage.SelectCharacters()

			if (err != nil) != tt.expectError {
				t.Errorf("SelectCharacters() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalCharactersSlice(characters, tt.expectedCharacters) {
				t.Errorf("SelectCharacters() результат = %v, ожидаемый %v", characters, tt.expectedCharacters)
			}
		})
	}
}

func TestUpdateCharacter(t *testing.T) {
	tests := []struct {
		name              string
		id                int64
		newCharacter      models.Character
		expectedCharacter models.Character
		expectError       bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newCharacter: models.Character{
				Name:  "Обновленный персонаж",
				Slug:  "updated-character",
				Color: "purple",
				Emotions: map[int64]int64{
					1: 1,
					2: 2,
				},
			},
			expectedCharacter: models.Character{
				Id:    1,
				Name:  "Обновленный персонаж",
				Slug:  "updated-character",
				Color: "purple",
				Emotions: map[int64]int64{
					1: 1,
					2: 2,
				},
			},
			expectError: false,
		},
		{
			name: "Персонаж не найден",
			id:   999,
			newCharacter: models.Character{
				Name: "Новый персонаж",
			},
			expectedCharacter: models.Character{},
			expectError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockCharacterStorage{
				UpdateCharacterFunc: func(id int64, newCharacter models.Character) (models.Character, error) {
					if id == 1 {
						return tt.expectedCharacter, nil
					}
					return models.Character{}, errors.New("character data not update")
				},
			}

			character, err := mockStorage.UpdateCharacter(tt.id, tt.newCharacter)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateCharacter() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalCharacters(character, tt.expectedCharacter) {
				t.Errorf("UpdateCharacter() результат = %v, ожидаемый %v", character, tt.expectedCharacter)
			}
		})
	}
}

func TestDeleteCharacter(t *testing.T) {
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
			name:           "Персонаж не найден",
			id:             999,
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockCharacterStorage{
				DeleteCharacterFunc: func(id int64) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("character data not update")
				},
			}

			result, err := mockStorage.DeleteCharacter(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("DeleteCharacter() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeleteCharacter() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func equalCharacters(a, b models.Character) bool {
	return a.Id == b.Id &&
		a.Name == b.Name &&
		a.Slug == b.Slug &&
		a.Color == b.Color &&
		mapsEqual(a.Emotions, b.Emotions)
}

func mapsEqual(a, b map[int64]int64) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func equalCharactersSlice(a, b []models.Character) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equalCharacters(v, b[i]) {
			return false
		}
	}
	return true
}
