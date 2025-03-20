package character

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestCreateCharacter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock соединения: %v", err)
	}
	defer db.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("ошибка при открытии GORM подключения: %v", err)
	}

	expectedID := generateUniqueId()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO characters \(name,slug,color,emotions,id\) VALUES \(?,?,?,?,?\)`).
		WillReturnResult(sqlmock.NewResult(expectedID, 1))

	mock.ExpectCommit()

	id, err := CreateCharacter("Тестовый Персонаж", "test-character", gormDB)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if id != expectedID {
		t.Errorf("ожидаемый ID = %d, полученный ID = %d", expectedID, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %v\nожидаемый запрос: %s",
			err,
			`INSERT INTO characters ([^)]+) VALUES (\?,\?,\?,\?,\?)`)
	}
}

func TestCreateCharacter_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании mock соединения: %v", err)
	}
	defer db.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("ошибка при открытии GORM подключения: %v", err)
	}

	// Настраиваем ожидание начала транзакции
	mock.ExpectBegin()

	// Ожидаем INSERT с ошибкой
	mock.ExpectExec(`INSERT INTO characters \([^)]+\)`).
		WillReturnError(sql.ErrConnDone)

	// Ожидаем rollback транзакции
	mock.ExpectRollback()

	_, err = CreateCharacter("Тестовый Персонаж", "test-character", gormDB)

	// Проверяем наличие ошибки
	if err == nil {
		t.Errorf("ожидалась ошибка, но она не произошла")
	}

	// Проверяем выполнение всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %v", err)
	}
}
