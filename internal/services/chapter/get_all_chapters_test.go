package chapter

import (
	"db_novel_service/internal/models"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetChaptersByUserId(t *testing.T) {
	// Настройка моковой базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании моковой БД: %v", err)
	}
	defer db.Close()

	// Создание подключения к GORM через мок
	dialect := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		t.Fatalf("ошибка при создании подключения к БД: %v", err)
	}

	// Тестовые данные
	testUserId := int64(42)
	testChapter := models.Chapter{
		Id:         1,
		Name:       "Тестовая глава",
		StartNode:  1,
		Nodes:      []int64{1, 2, 3},
		Characters: []int64{1, 2},
		Status:     3, // опубликована
		UpdatedAt:  map[time.Time]int64{time.Now(): 1},
		Author:     testUserId,
	}

	// Тестовый админ
	testAdmin := models.Admin{
		Id:               testUserId,
		Name:             "Тестовый админ",
		Email:            "admin@test.com",
		Password:         "password",
		AdminStatus:      1, // сверхадмин
		CreatedChapters:  []int64{1, 2},
		RequestSent:      []int64{3, 4},
		RequestsReceived: []int64{5, 6},
	}

	// Тестовый игрок
	testPlayer := models.Player{
		Id:                testUserId,
		Name:              "Тестовый игрок",
		Email:             "player@test.com",
		Password:          "password",
		Admin:             false,
		CompletedChapters: []int64{1, 2},
		ChaptersProgress:  map[int64]int64{1: 1, 2: 2},
		SoundSettings:     100,
	}

	// Тестовые случаи
	tests := []struct {
		name           string
		userId         int64
		isPlayer       bool
		isAdmin        bool
		expectedStatus int
		wantErr        bool
	}{
		{
			name:           "Админ - успешно",
			userId:         testUserId,
			isPlayer:       false,
			isAdmin:        true,
			expectedStatus: 3,
			wantErr:        false,
		},
		{
			name:           "Игрок - успешно",
			userId:         testUserId,
			isPlayer:       true,
			isAdmin:        false,
			expectedStatus: 3,
			wantErr:        false,
		},
		{
			name:           "Нет пользователя",
			userId:         testUserId,
			isPlayer:       false,
			isAdmin:        false,
			expectedStatus: 0,
			wantErr:        true,
		},
	}

	// Запускаем каждый тестовый случай
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Очищаем все ожидания
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			// Настройка ожидаемых запросов
			if tt.isAdmin {
				// Ожидаем успешный поиск админа с JSON полями
				createdChaptersJSON, _ := json.Marshal(testAdmin.CreatedChapters)
				requestSentJSON, _ := json.Marshal(testAdmin.RequestSent)
				requestsReceivedJSON, _ := json.Marshal(testAdmin.RequestsReceived)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, password, admin_status, COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw, COALESCE(request_sent::TEXT, '[]') as request_sent_raw, COALESCE(requests_received::TEXT, '[]') as requests_received_raw FROM admins WHERE id = $1 LIMIT 1`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "admin_status", "created_chapters_raw", "request_sent_raw", "requests_received_raw"}).
						AddRow(testUserId, testAdmin.Name, testAdmin.Email, testAdmin.Password, testAdmin.AdminStatus,
							string(createdChaptersJSON), string(requestSentJSON), string(requestsReceivedJSON)))

				// Ожидаем получение глав для админа
				nodesJSON, _ := json.Marshal(testChapter.Nodes)
				charactersJSON, _ := json.Marshal(testChapter.Characters)
				updatedAtJSON, _ := json.Marshal(testChapter.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, start_node, CAST(COALESCE(nodes, '[]'::json) AS TEXT) as nodes_raw, CAST(COALESCE(characters, '[]'::json) AS TEXT) as characters_raw, status, CAST(updated_at AS TEXT) as updated_at_raw, author FROM chapters`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "start_node", "nodes_raw", "characters_raw", "status", "updated_at_raw", "author"}).
						AddRow(testChapter.Id, testChapter.Name, testChapter.StartNode,
							string(nodesJSON), string(charactersJSON),
							testChapter.Status, string(updatedAtJSON), testChapter.Author))
			} else if tt.isPlayer {
				// Ожидаем успешный поиск игрока с JSON полями
				completedChaptersJSON, _ := json.Marshal(testPlayer.CompletedChapters)
				chaptersProgressJSON, _ := json.Marshal(testPlayer.ChaptersProgress)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, email, phone, password, admin, CAST(completed_chapters AS TEXT) as completed_chapters_raw, CAST(chapters_progress AS TEXT) as chapters_progress_raw, sound_settings FROM players WHERE id = $1 LIMIT 1`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "phone", "password", "admin", "completed_chapters_raw", "chapters_progress_raw", "sound_settings"}).
						AddRow(testUserId, testPlayer.Name, testPlayer.Email, testPlayer.Phone, testPlayer.Password,
							testPlayer.Admin, string(completedChaptersJSON), string(chaptersProgressJSON), testPlayer.SoundSettings))

				// Ожидаем получение опубликованных глав
				nodesJSON, _ := json.Marshal(testChapter.Nodes)
				charactersJSON, _ := json.Marshal(testChapter.Characters)
				updatedAtJSON, _ := json.Marshal(testChapter.UpdatedAt)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, start_node, CAST(nodes AS TEXT) as nodes_raw, CAST(characters AS TEXT) as characters_raw, status, CAST(updated_at AS TEXT) as updated_at_raw, author FROM chapters WHERE status = $1`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "start_node", "nodes_raw", "characters_raw", "status", "updated_at_raw", "author"}).
						AddRow(testChapter.Id, testChapter.Name, testChapter.StartNode,
							string(nodesJSON), string(charactersJSON),
							testChapter.Status, string(updatedAtJSON), testChapter.Author))
			}

			// Выполняем тестируемую функцию
			chapters, err := GetChaptersByUserId(gormDB, tt.userId)

			// Проверяем результаты
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChaptersByUserId() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(chapters) != 1 {
					t.Errorf("GetChaptersByUserId() количество глав = %d, хотим 1", len(chapters))
				}
				if chapters[0].Status != tt.expectedStatus {
					t.Errorf("GetChaptersByUserId() статус главы = %d, хотим %d", chapters[0].Status, tt.expectedStatus)
				}
			}

			// Проверяем, что все ожидаемые запросы были выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("не все ожидания были выполнены: %s", err)
			}
		})
	}
}
