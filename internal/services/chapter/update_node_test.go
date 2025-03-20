package chapter

import (
	"db_novel_service/internal/models"
	"encoding/json"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUpdateChapter(t *testing.T) {
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

	gormDB, err := gorm.Open(dialect, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		t.Fatalf("ошибка при создании подключения к БД: %v", err)
	}

	// Тестовые данные
	testChapterId := int64(42)
	testChapter := models.Chapter{
		Id:         testChapterId,
		Name:       "Тестовая глава",
		StartNode:  1,
		Nodes:      []int64{1, 2, 3},
		Characters: []int64{1, 2},
		Status:     3,
		UpdatedAt:  map[time.Time]int64{time.Now(): 1},
		Author:     1,
	}

	// Тестовый случай
	tests := []struct {
		label             string
		id                int64
		name              string
		nodes             []int64
		characters        []int64
		updateAuthorId    int64
		startNode         int64
		status            int
		wantErr           bool
		expectedUpdatedAt map[time.Time]int64
	}{
		{
			label:             "Успешное обновление всех полей",
			id:                testChapterId,
			name:              "Новое название главы",
			nodes:             []int64{4, 5, 6},
			characters:        []int64{3, 4},
			updateAuthorId:    2,
			startNode:         2,
			status:            2,
			wantErr:           false,
			expectedUpdatedAt: map[time.Time]int64{time.Now(): 2},
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
			nodesJSON, err := json.Marshal(tt.nodes)
			if err != nil {
				t.Fatal(err)
			}
			charactersJSON, err := json.Marshal(tt.characters)
			if err != nil {
				t.Fatal(err)
			}

			updatedAt := map[time.Time]int64{
				time.Now(): tt.updateAuthorId,
			}
			updatedAtJSON, err := json.Marshal(updatedAt)
			if err != nil {
				t.Fatal(err)
			}

			// Включаем логирование для отладки
			gormDB = gormDB.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})

			// Ожидаем успешный поиск главы
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, start_node, CAST(nodes AS TEXT) as nodes_raw, CAST(characters AS TEXT) as characters_raw, status, CAST(updated_at AS TEXT) as updated_at_raw, author FROM chapters WHERE id = $1 LIMIT 1`)).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "start_node", "nodes_raw", "characters_raw", "status", "updated_at_raw", "author"}).
					AddRow(testChapter.Id, testChapter.Name, testChapter.StartNode,
						string(nodesJSON), string(charactersJSON),
						testChapter.Status, string(updatedAtJSON), testChapter.Author))

			// Ожидаем обновление записи
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "chapters" SET "characters"=$1,"nodes"=$2,"start_node"=$3,"status"=$4,"updated_at"=$5,"name"=$6 WHERE id = $7`)).
				WithArgs(
					json.RawMessage(string(charactersJSON)),
					json.RawMessage(string(nodesJSON)),
					tt.startNode,
					tt.status,
					sqlmock.AnyArg(),
					tt.name,
					testChapterId,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			// Выполняем тестируемую функцию
			err = UpdateChapter(
				tt.id,
				tt.name,
				tt.nodes,
				tt.characters,
				tt.updateAuthorId,
				tt.startNode,
				tt.status,
				gormDB,
			)

			// Проверяем результаты
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateChapter() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Проверяем, что все ожидаемые запросы были выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("не все ожидания были выполнены: %s", err)
			}
		})
	}
}
