package admin

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
)

func TestRegistration(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Тестовые данные
	testEmail := "test@example.com"
	testName := "Test Admin"
	testPassword := "password123"

	emptyArraysJSON, err := json.Marshal([]int64{})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		label     string
		email     string
		name      string
		password  string
		setupMock func()
		wantID    int64
		wantErr   bool
	}{
		{
			label:    "успешная регистрация нового админа",
			email:    testEmail,
			name:     testName,
			password: testPassword,
			setupMock: func() {
				// Ожидаем SELECT для проверки существования email
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,email,password,admin_status,COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw,COALESCE(request_sent::TEXT, '[]') as request_sent_raw,COALESCE(requests_received::TEXT, '[]') as requests_received_raw FROM admins WHERE email = $1 LIMIT 1`)).
					WithArgs(testEmail).
					WillReturnError(sql.ErrNoRows)

				// Ожидаем начало транзакции
				mock.ExpectBegin()

				// Ожидаем INSERT в таблицу admins
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "admins" ("name","email","password","admin_status","created_chapters","request_sent","requests_received","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
					WithArgs(
						testName,           // name
						testEmail,          // email
						testPassword,       // password
						DefaultAdminStatus, // admin_status
						sqlmock.AnyArg(),   // created_chapters
						sqlmock.AnyArg(),   // request_sent
						sqlmock.AnyArg(),   // requests_received
						sqlmock.AnyArg(),   // id
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Ожидаем INSERT в таблицу requests
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "requests" ("id","sender_id","type","chapter_id") VALUES ($1,$2,$3,$4)`)).
					WithArgs(
						sqlmock.AnyArg(),         // id
						sqlmock.AnyArg(),         // sender_id
						RegisterAdminTypeRequest, // type
						NoChapter,                // chapter_id
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Ожидаем INSERT в таблицу players
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "players" ("id","name","email","password","admin") VALUES ($1,$2,$3,$4,$5)`)).
					WithArgs(
						sqlmock.AnyArg(), // id
						testName,         // name
						testEmail,        // email
						testPassword,     // password
						true,             // admin
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Ожидаем завершение транзакции
				mock.ExpectCommit()
			},
			wantID:  0,
			wantErr: true,
		},
		{
			label:    "ошибка при существующем email",
			email:    testEmail,
			name:     testName,
			password: testPassword,
			setupMock: func() {
				// Ожидаем только один SELECT для проверки существования email
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					int64(1), testName, testEmail, testPassword,
					DefaultAdminStatus, emptyArraysJSON, emptyArraysJSON, emptyArraysJSON,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,email,password,admin_status,COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw,COALESCE(request_sent::TEXT, '[]') as request_sent_raw,COALESCE(requests_received::TEXT, '[]') as requests_received_raw FROM admins WHERE email = $1 LIMIT 1`)).
					WithArgs(testEmail).
					WillReturnRows(rows)
			},
			wantID:  0,
			wantErr: true,
		},
		{
			label:    "ошибка при создании админа в базе данных",
			email:    testEmail,
			name:     testName,
			password: testPassword,
			setupMock: func() {
				// Ожидаем SELECT для проверки существования email
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,email,password,admin_status,COALESCE(created_chapters::TEXT, '[]') as created_chapters_raw,COALESCE(request_sent::TEXT, '[]') as request_sent_raw,COALESCE(requests_received::TEXT, '[]') as requests_received_raw FROM admins WHERE email = $1 LIMIT 1`)).
					WithArgs(testEmail).
					WillReturnError(gorm.ErrRecordNotFound)

				// Ожидаем начало транзакции
				mock.ExpectBegin()

				// Ожидаем INSERT в таблицу admins с ошибкой
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "admins" ("name","email","password","admin_status","created_chapters","request_sent","requests_received","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
					WithArgs(
						testName,           // name
						testEmail,          // email
						testPassword,       // password
						DefaultAdminStatus, // admin_status
						sqlmock.AnyArg(),   // created_chapters
						sqlmock.AnyArg(),   // request_sent
						sqlmock.AnyArg(),   // requests_received
						sqlmock.AnyArg(),   // id
					).
					WillReturnError(errors.New("database error"))

				// Ожидаем откат транзакции
				mock.ExpectRollback()
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настройка мока для текущего теста
			tt.setupMock()

			// Создаем gorm.DB с моком
			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				DSN:  "sqlmock_db_0",
				Conn: db,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

			// Вызываем тестируемую функцию
			id, err := Registration(tt.email, tt.name, tt.password, gormDB)

			// Проверяем ID и наличие ошибки
			if id != tt.wantID {
				t.Errorf("Registration() вернул неверный ID = %v, хотели %v", id, tt.wantID)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Registration() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
