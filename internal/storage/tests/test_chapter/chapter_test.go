package test_chapter

import (
	"db_novel_service/internal/models"
	"errors"
	"testing"
	"time"
)

func TestRegisterChapter(t *testing.T) {
	tests := []struct {
		name           string
		inputChapter   models.Chapter
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputChapter: models.Chapter{
				Name:       "Тестовая глава",
				Nodes:      []int64{1, 2, 3},
				Characters: []int64{4, 5},
				UpdatedAt: map[time.Time]int64{
					time.Now(): 1,
				},
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Ошибка валидации",
			inputChapter: models.Chapter{
				Name: "", // пустое имя
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				RegisterChapterFunc: func(chapter models.Chapter) (int64, error) {
					if chapter.Name == "" {
						return 0, errors.New("не удалось создать запись главы")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterChapter(tt.inputChapter)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterChapter() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterChapter() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectChapterWithId(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		expectedChapter models.Chapter
		expectError     bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedChapter: models.Chapter{
				Id:         1,
				Name:       "Тестовая глава",
				Nodes:      []int64{1, 2, 3},
				Characters: []int64{4, 5},
				UpdatedAt: map[time.Time]int64{
					time.Now(): 1,
				},
			},
			expectError: false,
		},
		{
			name:            "Глава не найдена",
			id:              999,
			expectedChapter: models.Chapter{},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				SelectChapterWithIdFunc: func(id int64) (models.Chapter, error) {
					if id == 1 {
						return tt.expectedChapter, nil
					}
					return models.Chapter{}, errors.New("chapter data not found")
				},
			}

			chapter, err := mockStorage.SelectChapterWithId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectChapterWithId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalChapters(chapter, tt.expectedChapter) {
				t.Errorf("SelectChapterWithId() результат = %v, ожидаемый %v", chapter, tt.expectedChapter)
			}
		})
	}
}

func TestUpdateChapter(t *testing.T) {
	tests := []struct {
		name            string
		id              int64
		newChapter      models.Chapter
		expectedChapter models.Chapter
		expectError     bool
	}{
		{
			name: "Успешное обновление",
			id:   1,
			newChapter: models.Chapter{
				Name:       "Обновленная глава",
				Nodes:      []int64{10, 20, 30},
				Characters: []int64{40, 50},
				UpdatedAt: map[time.Time]int64{
					time.Now(): 1,
				},
			},
			expectedChapter: models.Chapter{
				Id:         1,
				Name:       "Обновленная глава",
				Nodes:      []int64{10, 20, 30},
				Characters: []int64{40, 50},
				UpdatedAt: map[time.Time]int64{
					time.Now(): 1,
				},
			},
			expectError: false,
		},
		{
			name: "Глава не найдена",
			id:   999,
			newChapter: models.Chapter{
				Name: "Новая глава",
			},
			expectedChapter: models.Chapter{},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				UpdateChapterFunc: func(id int64, newChapter models.Chapter) (models.Chapter, error) {
					if id == 1 {
						return tt.expectedChapter, nil
					}
					return models.Chapter{}, errors.New("chapter data not update")
				},
			}

			chapter, err := mockStorage.UpdateChapter(tt.id, tt.newChapter)

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateChapter() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalChapters(chapter, tt.expectedChapter) {
				t.Errorf("UpdateChapter() результат = %v, ожидаемый %v", chapter, tt.expectedChapter)
			}
		})
	}
}

func TestDeleteChapter(t *testing.T) {
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
			name:           "Глава не найдена",
			id:             999,
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				DeleteChapterFunc: func(id int64) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("payment data not update")
				},
			}

			result, err := mockStorage.DeleteChapter(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("DeleteChapter() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeleteChapter() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestGetChaptersForAdmin(t *testing.T) {
	tests := []struct {
		name             string
		expectedChapters []models.Chapter
		expectError      bool
	}{
		{
			name: "Успешный поиск",
			expectedChapters: []models.Chapter{
				{
					Id:         1,
					Name:       "Глава 1",
					Nodes:      []int64{1, 2, 3},
					Characters: []int64{4, 5},
					UpdatedAt: map[time.Time]int64{
						time.Now(): 1,
					},
				},
				{
					Id:         2,
					Name:       "Глава 2",
					Nodes:      []int64{6, 7, 8},
					Characters: []int64{9, 10},
					UpdatedAt: map[time.Time]int64{
						time.Now(): 2,
					},
				},
			},
			expectError: false,
		},
		{
			name:             "Ошибка получения глав",
			expectedChapters: []models.Chapter{},
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				GetChaptersForAdminFunc: func() ([]models.Chapter, error) {
					if len(tt.expectedChapters) > 0 {
						return tt.expectedChapters, nil
					}
					return nil, errors.New("failed to fetch published chapters")
				},
			}

			chapters, err := mockStorage.GetChaptersForAdmin()

			if (err != nil) != tt.expectError {
				t.Errorf("GetChaptersForAdmin() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalChaptersSlice(chapters, tt.expectedChapters) {
				t.Errorf("GetChaptersForAdmin() результат = %v, ожидаемый %v", chapters, tt.expectedChapters)
			}
		})
	}
}

func TestFindPublishedChapters(t *testing.T) {
	tests := []struct {
		name             string
		expectedChapters []models.Chapter
		expectError      bool
	}{
		{
			name: "Успешный поиск",
			expectedChapters: []models.Chapter{
				{
					Id:         1,
					Name:       "Опубликованная глава",
					Status:     3,
					Nodes:      []int64{1, 2, 3},
					Characters: []int64{4, 5},
					UpdatedAt: map[time.Time]int64{
						time.Now(): 1,
					},
				},
			},
			expectError: false,
		},
		{
			name:             "Главы не найдены",
			expectedChapters: []models.Chapter{},
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockChapterStorage{
				FindPublishedChaptersFunc: func() ([]models.Chapter, error) {
					if len(tt.expectedChapters) > 0 {
						return tt.expectedChapters, nil
					}
					return nil, errors.New("failed to fetch published chapters")
				},
			}

			chapters, err := mockStorage.FindPublishedChapters()

			if (err != nil) != tt.expectError {
				t.Errorf("FindPublishedChapters() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalChaptersSlice(chapters, tt.expectedChapters) {
				t.Errorf("FindPublishedChapters() результат = %v, ожидаемый %v", chapters, tt.expectedChapters)
			}
		})
	}
}

func equalChapters(a, b models.Chapter) bool {
	return a.Id == b.Id &&
		a.Name == b.Name &&
		jsonEqual(a.Nodes, b.Nodes) &&
		jsonEqual(a.Characters, b.Characters) &&
		mapsEqual(a.UpdatedAt, b.UpdatedAt)
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

func mapsEqual(a, b map[time.Time]int64) bool {
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

func equalChaptersSlice(a, b []models.Chapter) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equalChapters(v, b[i]) {
			return false
		}
	}
	return true
}
