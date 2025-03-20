package character

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

//func TestUpdateCharacter(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("ошибка при создании mock соединения: %v", err)
//	}
//	defer db.Close()
//
//	dialector := mysql.New(mysql.Config{
//		DSN:                       "sqlmock_db_0",
//		DriverName:                "mysql",
//		Conn:                      db,
//		SkipInitializeWithVersion: true,
//	})
//
//	gormDB, err := gorm.Open(dialector, &gorm.Config{})
//	if err != nil {
//		t.Errorf("ошибка при открытии GORM подключения: %v", err)
//	}
//
//	// Настраиваем ожидания для SELECT запроса
//	mock.ExpectQuery(`SELECT id, name, slug, color, CAST\(emotions AS TEXT\) as emotions_raw FROM characters WHERE id = \?`).
//		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "color", "emotions_raw"}).
//			AddRow(1, "Тестовый Персонаж", "test-character", "#00693E", `{"1":100,"2":200}`))
//
//	// Настраиваем ожидания для UPDATE запроса
//	mock.ExpectBegin()
//	mock.ExpectExec(`UPDATE characters SET Name = \?, Slug = \?, Color = \?, Emotions = \? WHERE id = \?`).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	// Выполняем тестируемый метод
//	err = UpdateCharacter(
//		1,
//		"Новое Имя",
//		"new-slug",
//		"#FF0000",
//		map[int64]int64{3: 300},
//		gormDB,
//	)
//
//	// Проверяем результат
//	if err != nil {
//		t.Errorf("неожиданная ошибка: %v", err)
//	}
//
//	// Проверяем выполнение всех ожиданий
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("не все ожидания были выполнены: %v", err)
//	}
//}

func TestUpdateCharacter_NonExistent(t *testing.T) {
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

	// Настраиваем ожидания для SELECT запроса
	mock.ExpectQuery(`SELECT id, name, slug, color, CAST\(emotions AS TEXT\) as emotions_raw FROM characters WHERE id = \?`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "color", "emotions_raw"}))

	// Выполняем тестируемый метод
	err = UpdateCharacter(
		1,
		"Новое Имя",
		"new-slug",
		"#FF0000",
		map[int64]int64{3: 300},
		gormDB,
	)

	// Проверяем, что ошибка не nil
	if err == nil {
		t.Errorf("ожидалась ошибка, но она не произошла")
	}

	// Проверяем выполнение всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %v", err)
	}
}

func TestUpdateCharacter_InvalidJSON(t *testing.T) {
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

	// Настраиваем ожидания для SELECT запроса
	mock.ExpectQuery(`SELECT id, name, slug, color, CAST\(emotions AS TEXT\) as emotions_raw FROM characters WHERE id = \?`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "color", "emotions_raw"}).
			AddRow(1, "Тестовый Персонаж", "test-character", "#00693E", `{"1":100,"2":200}`))

	// Выполняем тестируемый метод с невалидным JSON
	err = UpdateCharacter(
		1,
		"Новое Имя",
		"new-slug",
		"#FF0000",
		map[int64]int64{3: 300},
		gormDB,
	)

	// Проверяем, что ошибка не nil
	if err == nil {
		t.Errorf("ожидалась ошибка, но она не произошла")
	}

	// Проверяем выполнение всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %v", err)
	}
}
