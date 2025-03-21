package character

import (
	"gorm.io/gorm"
	"math/rand"
	"time"
	"vn/internal/models"
	"vn/internal/storage"
)

func CreateCharacter(name string, slug string, db *gorm.DB) (int64, error) {

	id := generateUniqueId()

	newCharacter := models.Character{
		Id:       id,
		Name:     name,
		Slug:     slug,
		Emotions: map[int64]int64{},
		Color:    "#00693E",
	}

	_, err := storage.RegisterCharacter(db, newCharacter)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func generateUniqueId() int64 {
	// Получаем текущее время в миллисекундах (48 бит)
	timestamp := time.Now().UnixMilli()

	// Генерируем 16 случайных бит
	random := rand.Int31n(1 << 16)

	// Объединяем timestamp и random в 64-битное число
	return (int64(timestamp) << 16) | int64(random)
}
