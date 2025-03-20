package character

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGetCharacters(t *testing.T) {
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

	// Создаем тестовые данные
	expectedJSON := `{"1": 100, "2": 200}`
	mock.ExpectQuery(`SELECT id, name, slug, color, CAST\(emotions AS TEXT\) as emotions_raw FROM characters`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "color", "emotions_raw"}).
			AddRow(1, "Тестовый Персонаж", "test-character", "#00693E", expectedJSON))

	// Выполняем тестируемый метод
	characters, err := GetCharacters(gormDB)

	// Проверяем отсутствие ошибок
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	// Проверяем, что данные корректно десериализованы
	if characters == nil {
		t.Errorf("результат должен быть не nil")
	}

	if len(*characters) != 1 {
		t.Errorf("ожидается 1 персонаж, получено: %d", len(*characters))
	}

	// Проверяем корректность данных
	c := (*characters)[0]
	if c.Id != 1 {
		t.Errorf("ожидаемый ID = 1, полученный ID = %d", c.Id)
	}

	if c.Name != "Тестовый Персонаж" {
		t.Errorf("ожидаемое имя = 'Тестовый Персонаж', полученное имя = '%s'", c.Name)
	}

	if c.Slug != "test-character" {
		t.Errorf("ожидаемый slug = 'test-character', полученный slug = '%s'", c.Slug)
	}

	if c.Color != "#00693E" {
		t.Errorf("ожидаемый цвет = '#00693E', полученный цвет = '%s'", c.Color)
	}

	// Проверяем корректность JSON
	if len(c.Emotions) != 2 {
		t.Errorf("ожидается 2 эмоции, получено: %d", len(c.Emotions))
	}

	if c.Emotions[1] != 100 {
		t.Errorf("ожидаемое значение для эмоции 1 = 100, полученное = %d", c.Emotions[1])
	}

	if c.Emotions[2] != 200 {
		t.Errorf("ожидаемое значение для эмоции 2 = 200, полученное = %d", c.Emotions[2])
	}

	// Проверяем выполнение всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %v", err)
	}
}
