package admin

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"vn/internal/models"
)

func TestChangeAdmin(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	// Создаем тестового админа с всеми полями
	testAdmin := models.Admin{
		Id:               1,
		Name:             "Test Admin",
		Email:            "test@example.com",
		Password:         "password123",
		AdminStatus:      1,
		CreatedChapters:  []int64{1, 2, 3},
		RequestSent:      []int64{4, 5},
		RequestsReceived: []int64{6, 7},
	}

	// Преобразуем массивы в JSON строки
	createdChaptersJSON, err := json.Marshal(testAdmin.CreatedChapters)
	if err != nil {
		t.Fatal(err)
	}
	requestSentJSON, err := json.Marshal(testAdmin.RequestSent)
	if err != nil {
		t.Fatal(err)
	}
	requestsReceivedJSON, err := json.Marshal(testAdmin.RequestsReceived)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		newName            string
		newEmail           string
		newPassword        string
		newAdminStatus     int
		newCreatedChapters []int64
		setupMock          func()
		wantErr            bool
	}{
		{
			name:               "успешное обновление всех полей",
			newName:            "Updated Admin",
			newEmail:           "updated@example.com",
			newPassword:        "newpass123",
			newAdminStatus:     2,
			newCreatedChapters: []int64{4, 5, 6},
			setupMock: func() {
				// Ожидаем SELECT для получения админа
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					testAdmin.Id, testAdmin.Name, testAdmin.Email, testAdmin.Password,
					testAdmin.AdminStatus, createdChaptersJSON, requestSentJSON, requestsReceivedJSON,
				)
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE id = ?").
					WithArgs(testAdmin.Id).
					WillReturnRows(rows)

				// Ожидаем UPDATE с правильным количеством аргументов
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"admins\"").
					WithArgs(
						sqlmock.AnyArg(), // admin_status
						sqlmock.AnyArg(), // created_chapters
						sqlmock.AnyArg(), // email
						sqlmock.AnyArg(), // name
						sqlmock.AnyArg(), // password
						sqlmock.AnyArg(), // request_sent
						sqlmock.AnyArg(), // requests_received
						testAdmin.Id,     // id
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name:    "ошибка при получении админа",
			newName: "Updated Admin",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE id = ?").
					WithArgs(testAdmin.Id).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
		{
			name:    "ошибка при обновлении",
			newName: "Updated Admin",
			setupMock: func() {
				// Успешно получаем админа
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					testAdmin.Id, testAdmin.Name, testAdmin.Email, testAdmin.Password,
					testAdmin.AdminStatus, createdChaptersJSON, requestSentJSON, requestsReceivedJSON,
				)
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE id = ?").
					WithArgs(testAdmin.Id).
					WillReturnRows(rows)

				// Ожидаем ошибку при UPDATE
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE \"admins\"").
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
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
			}), &gorm.Config{})
			if err != nil {
				t.Fatal(err)
			}

			// Вызываем тестируемую функцию
			err = ChangeAdmin(
				testAdmin.Id,
				tt.newName,
				tt.newEmail,
				tt.newPassword,
				tt.newAdminStatus,
				tt.newCreatedChapters,
				gormDB,
			)

			// Проверяем наличие ошибки
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeAdmin() ошибка = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Проверяем, что все ожидания были выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("не все ожидания были выполнены: %v", err)
			}
		})
	}
}
