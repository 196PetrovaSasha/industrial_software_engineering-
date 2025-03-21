package test_player

import (
	"errors"
	"fmt"
	"testing"
	"vn/internal/models"
)

func TestRegisterPlayer(t *testing.T) {
	tests := []struct {
		name           string
		inputPlayer    models.Player
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputPlayer: models.Player{
				Name:              "Тестовый игрок",
				Email:             "test@example.com",
				Phone:             "+79999999999",
				Password:          "password123",
				Admin:             false,
				CompletedChapters: []int64{},
				ChaptersProgress:  map[int64]int64{},
				SoundSettings:     100,
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Пустой email",
			inputPlayer: models.Player{
				Name:  "Тестовый игрок",
				Email: "", // пустой email
			},
			expectedResult: 0,
			expectedErr:    true,
		},
		{
			name: "Неверный формат JSON",
			inputPlayer: models.Player{
				Name:              "Тестовый игрок",
				Email:             "test@example.com",
				CompletedChapters: []int64{1},
				ChaptersProgress:  map[int64]int64{1: 1},
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockPlayerStorage{
				RegisterPlayerFunc: func(player models.Player) (int64, error) {
					if player.Email == "" {
						return 0, errors.New("email не может быть пустым")
					}

					if len(player.CompletedChapters) > 0 || len(player.ChaptersProgress) > 0 {
						return 0, fmt.Errorf("неверный формат JSON для полей CompletedChapters или ChaptersProgress")
					}

					return 1, nil
				},
			}

			result, err := mockStorage.RegisterPlayer(tt.inputPlayer)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterPlayer() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterPlayer() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectPlayerWIthEmail(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		expectedPlayer models.Player
		expectError    bool
	}{
		{
			name:  "Успешный поиск",
			email: "test@example.com",
			expectedPlayer: models.Player{
				Id:                1,
				Name:              "Тестовый игрок",
				Email:             "test@example.com",
				Phone:             "+79999999999",
				Password:          "password123",
				Admin:             false,
				CompletedChapters: []int64{},
				ChaptersProgress:  map[int64]int64{},
				SoundSettings:     100,
			},
			expectError: false,
		},
		{
			name:           "Игрок не найден",
			email:          "nonexistent@example.com",
			expectedPlayer: models.Player{},
			expectError:    true,
		},
		{
			name:           "Ошибка десериализации JSON",
			email:          "json-error@example.com",
			expectedPlayer: models.Player{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockPlayerStorage{
				SelectPlayerWIthEmailFunc: func(email string) (models.Player, error) {
					switch email {
					case "test@example.com":
						return tt.expectedPlayer, nil
					case "nonexistent@example.com":
						return models.Player{}, errors.New("player data not found")
					case "json-error@example.com":
						return models.Player{}, fmt.Errorf("failed to unmarshal completed_chapters: invalid character")
					default:
						return models.Player{}, errors.New("unexpected email")
					}
				},
			}

			player, err := mockStorage.SelectPlayerWIthEmail(tt.email)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectPlayerWIthEmail() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalPlayers(player, tt.expectedPlayer) {
				t.Errorf("SelectPlayerWIthEmail() результат = %v, ожидаемый %v", player, tt.expectedPlayer)
			}
		})
	}
}

func TestSelectPlayerWIthId(t *testing.T) {
	tests := []struct {
		name           string
		id             int64
		expectedPlayer models.Player
		expectError    bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedPlayer: models.Player{
				Id:                1,
				Name:              "Тестовый игрок",
				Email:             "test@example.com",
				Phone:             "+79999999999",
				Password:          "password123",
				Admin:             false,
				CompletedChapters: []int64{},
				ChaptersProgress:  map[int64]int64{},
				SoundSettings:     100,
			},
			expectError: false,
		},
		{
			name:           "Игрок не найден",
			id:             999,
			expectedPlayer: models.Player{},
			expectError:    true,
		},
		{
			name:           "Ошибка десериализации JSON",
			id:             888,
			expectedPlayer: models.Player{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockPlayerStorage{
				SelectPlayerWIthIdFunc: func(id int64) (models.Player, error) {
					switch id {
					case 1:
						return tt.expectedPlayer, nil
					case 999:
						return models.Player{}, errors.New("player data not found")
					case 888:
						return models.Player{}, fmt.Errorf("failed to unmarshal completed_chapters: invalid character")
					default:
						return models.Player{}, errors.New("unexpected id")
					}
				},
			}

			player, err := mockStorage.SelectPlayerWIthId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectPlayerWIthId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalPlayers(player, tt.expectedPlayer) {
				t.Errorf("SelectPlayerWIthId() результат = %v, ожидаемый %v", player, tt.expectedPlayer)
			}
		})
	}
}

func TestUpdatePlayer(t *testing.T) {
	tests := []struct {
		name           string
		id             int64
		newPlayer      models.Player
		expectedPlayer models.Player
		expectError    bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newPlayer: models.Player{
				Name:              "Обновленный игрок",
				Email:             "updated@example.com",
				Phone:             "+79999999999",
				Password:          "newpassword123",
				Admin:             false,
				CompletedChapters: []int64{},
				ChaptersProgress:  map[int64]int64{},
				SoundSettings:     100,
			},
			expectedPlayer: models.Player{
				Id:                1,
				Name:              "Обновленный игрок",
				Email:             "updated@example.com",
				Phone:             "+79999999999",
				Password:          "newpassword123",
				Admin:             false,
				CompletedChapters: []int64{},
				ChaptersProgress:  map[int64]int64{},
				SoundSettings:     100,
			},
			expectError: false,
		},
		{
			name: "Игрок не найден",
			id:   999,
			newPlayer: models.Player{
				Name: "Новый игрок",
			},
			expectedPlayer: models.Player{},
			expectError:    true,
		},
		{
			name: "Ошибка JSON",
			id:   888,
			newPlayer: models.Player{
				CompletedChapters: []int64{1},
				ChaptersProgress:  map[int64]int64{1: 1},
			},
			expectedPlayer: models.Player{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockPlayerStorage{
				UpdatePlayerFunc: func(id int64, newPlayer models.Player) (models.Player, error) {
					switch id {
					case 1:
						if len(newPlayer.CompletedChapters) > 0 || len(newPlayer.ChaptersProgress) > 0 {
							return models.Player{}, fmt.Errorf("неверный формат JSON для полей CompletedChapters или ChaptersProgress")
						}
						return tt.expectedPlayer, nil
					case 999:
						return models.Player{}, errors.New("player data not updated")
					case 888:
						return models.Player{}, fmt.Errorf("ошибка маршалинга ChaptersProgress: invalid character")
					default:
						return models.Player{}, errors.New("unexpected id")
					}
				},
			}

			player, err := mockStorage.UpdatePlayer(tt.id, tt.newPlayer)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdatePlayer() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalPlayers(player, tt.expectedPlayer) {
				t.Errorf("UpdatePlayer() результат = %v, ожидаемый %v", player, tt.expectedPlayer)
			}
		})
	}
}

func TestDeletePlayer(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		expectedResult int64
		expectError    bool
	}{
		{
			name:           "Успешное удаление",
			email:          "test@example.com",
			expectedResult: 1,
			expectError:    false,
		},
		{
			name:           "Игрок не найден",
			email:          "nonexistent@example.com",
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockPlayerStorage{
				DeletePlayerFunc: func(email string) (int64, error) {
					switch email {
					case "test@example.com":
						return 1, nil
					case "nonexistent@example.com":
						return 0, errors.New("player data not update")
					default:
						return 0, errors.New("unexpected email")
					}
				},
			}

			result, err := mockStorage.DeletePlayer(tt.email)

			if (err != nil) != tt.expectError {
				t.Errorf("DeletePlayer() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeletePlayer() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func equalPlayers(a, b models.Player) bool {
	return a.Id == b.Id &&
		a.Name == b.Name &&
		a.Email == b.Email &&
		a.Phone == b.Phone &&
		a.Password == b.Password &&
		a.Admin == b.Admin &&
		slicesEqual(a.CompletedChapters, b.CompletedChapters) && // Changed to slicesEqual
		mapsEqual(a.ChaptersProgress, b.ChaptersProgress) &&
		a.SoundSettings == b.SoundSettings
}

// New function for comparing slices
func slicesEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Existing function for comparing maps
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

func mapsEqualSlice(a, b []int64) bool {
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
