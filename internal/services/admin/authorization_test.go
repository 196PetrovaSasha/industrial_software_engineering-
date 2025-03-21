package admin

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"vn/internal/models"
)

func TestAuthorization(t *testing.T) {
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
		name          string
		email         string
		password      string
		setupMock     func()
		wantErr       bool
		expectedAdmin *models.Admin
	}{
		{
			name:     "успешная авторизация",
			email:    testAdmin.Email,
			password: testAdmin.Password,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					testAdmin.Id, testAdmin.Name, testAdmin.Email, testAdmin.Password,
					testAdmin.AdminStatus, createdChaptersJSON, requestSentJSON, requestsReceivedJSON,
				)
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE email = ?").
					WithArgs(testAdmin.Email).
					WillReturnRows(rows)
			},
			wantErr:       false,
			expectedAdmin: &testAdmin,
		}, {
			name:     "неуспешная авторизация",
			email:    testAdmin.Email,
			password: "45678",
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					testAdmin.Id, testAdmin.Name, testAdmin.Email, testAdmin.Password,
					testAdmin.AdminStatus, createdChaptersJSON, requestSentJSON, requestsReceivedJSON,
				)
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE email = ?").
					WithArgs(testAdmin.Email).
					WillReturnRows(rows)
			},
			wantErr:       true,
			expectedAdmin: &testAdmin,
		},
		{
			name:     "неуспешная авторизация c email",
			email:    "testtttttt@example.com",
			password: testAdmin.Password,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "password", "admin_status",
					"created_chapters", "request_sent", "requests_received",
				}).AddRow(
					testAdmin.Id, testAdmin.Name, testAdmin.Email, testAdmin.Password,
					testAdmin.AdminStatus, createdChaptersJSON, requestSentJSON, requestsReceivedJSON,
				)
				mock.ExpectQuery("SELECT (.+) FROM admins WHERE email = ?").
					WithArgs("testtttttt@example.com"). // Исправлено: используем тестовый email вместо testAdmin.Email
					WillReturnRows(rows)
			},
			wantErr:       true,
			expectedAdmin: &testAdmin,
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
			admin, err := Authorization(tt.email, tt.password, gormDB)

			// Проверяем результат при успешной авторизации
			if !tt.wantErr && !reflect.DeepEqual(admin, tt.expectedAdmin) {
				t.Errorf("Authorization() admin = %v, expectedAdmin %v", admin, tt.expectedAdmin)
			}

			// Проверяем, что все ожидания были выполнены
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("не все ожидания были выполнены: %v", err)
			}
		})
	}
}
