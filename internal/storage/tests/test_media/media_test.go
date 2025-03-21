package test_media

import (
	"errors"
	"testing"
	"vn/internal/models"
)

func TestRegisterMedia(t *testing.T) {
	tests := []struct {
		name           string
		inputMedia     models.Media
		expectedResult int64
		expectedErr    bool
	}{
		{
			name: "Успешная регистрация",
			inputMedia: models.Media{
				FileData:    []byte{1, 2, 3},
				ContentType: "image/jpeg",
			},
			expectedResult: 1,
			expectedErr:    false,
		},
		{
			name: "Ошибка регистрации",
			inputMedia: models.Media{
				FileData: nil,
			},
			expectedResult: 0,
			expectedErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockMediaStorage{
				RegisterMediaFunc: func(media models.Media) (int64, error) {
					if media.FileData == nil {
						return 0, errors.New("media not created")
					}
					return 1, nil
				},
			}

			result, err := mockStorage.RegisterMedia(tt.inputMedia)

			if (err != nil) != tt.expectedErr {
				t.Errorf("RegisterMedia() ошибка = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("RegisterMedia() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelectMediaWIthId(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		expectedMedia models.Media
		expectError   bool
	}{
		{
			name: "Успешный поиск",
			id:   1,
			expectedMedia: models.Media{
				Id:          1,
				FileData:    []byte{1, 2, 3},
				ContentType: "image/jpeg",
			},
			expectError: false,
		},
		{
			name:          "Медиа не найдено",
			id:            999,
			expectedMedia: models.Media{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockMediaStorage{
				SelectMediaWIthIdFunc: func(id int64) (models.Media, error) {
					if id == 1 {
						return tt.expectedMedia, nil
					}
					return models.Media{}, errors.New("media data not found")
				},
			}

			media, err := mockStorage.SelectMediaWIthId(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("SelectMediaWIthId() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalMedia(media, tt.expectedMedia) {
				t.Errorf("SelectMediaWIthId() результат = %v, ожидаемый %v", media, tt.expectedMedia)
			}
		})
	}
}

func TestSelectMedia(t *testing.T) {
	tests := []struct {
		name          string
		expectedMedia []models.Media
		expectError   bool
	}{
		{
			name: "Успешный поиск",
			expectedMedia: []models.Media{
				{
					Id:          1,
					FileData:    []byte{1, 2, 3},
					ContentType: "image/jpeg",
				},
				{
					Id:          2,
					FileData:    []byte{4, 5, 6},
					ContentType: "image/png",
				},
			},
			expectError: false,
		},
		{
			name:          "Ошибка получения медиа",
			expectedMedia: []models.Media{},
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockMediaStorage{
				SelectMediaFunc: func() ([]models.Media, error) {
					if len(tt.expectedMedia) > 0 {
						return tt.expectedMedia, nil
					}
					return nil, errors.New("media not found")
				},
			}

			media, err := mockStorage.SelectMedia()

			if (err != nil) != tt.expectError {
				t.Errorf("SelectMedia() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && !equalMediaSlice(media, tt.expectedMedia) {
				t.Errorf("SelectMedia() результат = %v, ожидаемый %v", media, tt.expectedMedia)
			}
		})
	}
}

func TestDeleteMedia(t *testing.T) {
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
			name:           "Медиа не найдено",
			id:             999,
			expectedResult: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := &MockMediaStorage{
				DeleteMediaFunc: func(id int64) (int64, error) {
					if id == 1 {
						return 1, nil
					}
					return 0, errors.New("media data not update")
				},
			}

			result, err := mockStorage.DeleteMedia(tt.id)

			if (err != nil) != tt.expectError {
				t.Errorf("DeleteMedia() ошибка = %v, expectError %v", err, tt.expectError)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("DeleteMedia() результат = %v, ожидаемый %v", result, tt.expectedResult)
			}
		})
	}
}

func equalMedia(a, b models.Media) bool {
	return a.Id == b.Id &&
		string(a.FileData) == string(b.FileData) &&
		a.ContentType == b.ContentType
}

func equalMediaSlice(a, b []models.Media) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !equalMedia(v, b[i]) {
			return false
		}
	}
	return true
}
