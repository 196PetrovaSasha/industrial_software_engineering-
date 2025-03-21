package chapter

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestCreateDefaultChapter(t *testing.T) {
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
	testAuthorId := int64(42)

	// Тестовые случаи
	tests := []struct {
		name          string
		authorId      int64
		wantChapterId int64
		wantNodeId    int64
		wantErr       bool
	}{
		{
			name:          "Успешное создание главы",
			authorId:      testAuthorId,
			wantChapterId: 1,
			wantNodeId:    2,
			wantErr:       true,
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
			if !tt.wantErr {
				// Ожидаем создание главы с JSON полями
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "chapters" ("id","name","nodes","characters","status","author","start_node","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantChapterId))

				// Ожидаем обновление главы с временной меткой
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "chapters" SET "updated_at" = $1 WHERE "id" = $2`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantChapterId))

				// Ожидаем создание узла
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "nodes" ("id","slug","chapter_id") VALUES ($1,$2,$3) RETURNING "id"`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantNodeId))

				// Ожидаем обновление главы с начальным узлом
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "chapters" SET "start_node" = $1 WHERE "id" = $2`)).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.wantChapterId))
			} else {
				// Ожидаем ошибку при создании главы
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "chapters" ("id","name","nodes","characters","status","author","start_node","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
					WillReturnError(sql.ErrConnDone)
			}

			// Выполняем тестируемую функцию
			chapterId, nodeId, err := CreateDefaultChapter(tt.authorId, gormDB)

			// Проверяем результаты
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDefaultChapter() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if chapterId != tt.wantChapterId {
					t.Errorf("CreateDefaultChapter() chapterId = %d, хотим %d", chapterId, tt.wantChapterId)
				}
				if nodeId != tt.wantNodeId {
					t.Errorf("CreateDefaultChapter() nodeId = %d, хотим %d", nodeId, tt.wantNodeId)
				}
			}
		})
	}
}
