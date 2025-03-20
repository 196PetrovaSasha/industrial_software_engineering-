package admin

import (
	"db_novel_service/internal/models"
	"db_novel_service/internal/storage"
	"errors"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"time"
)

const (
	AdminNotFoundError = "admin data not found"

	DefaultAdminStatus = -1

	NoChapter                = -1
	RegisterAdminTypeRequest = 1
)

func Registration(email string, name string, password string, db *gorm.DB) (int64, error) {
	_, err := storage.SelectAdminWIthEmail(db, email)

	if err == nil {
		return 0, errors.New("admin with this email is already exist")
	}

	//if err.Error() != AdminNotFoundError {
	//	log.Println(err, "ошибка получения админа")
	//	return 0, err
	//}

	id := generateUniqueId()

	newAdmin := models.Admin{
		Id:               id,
		Email:            email,
		Password:         password,
		Name:             name,
		AdminStatus:      DefaultAdminStatus,
		CreatedChapters:  []int64{},
		RequestSent:      []int64{},
		RequestsReceived: []int64{},
	}

	_, err = storage.RegisterAdmin(db, newAdmin)

	if err != nil {
		return 0, err
	}

	_, err = CreateRequest(id, RegisterAdminTypeRequest, NoChapter, db)

	log.Println(err)

	if err != nil {
		return 0, errors.New("fail to send registration requests to another admin")
	}

	_, err = storage.RegisterPlayer(db, models.Player{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Admin:    true,
	})

	if err != nil {
		return 0, errors.New("error to create player for admin")
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

func CreateRequest(requestingAdminId int64, typeRequest int, requestedChapterId int64, db *gorm.DB) (int64, error) {

	id := generateUniqueId()

	log.Println(requestingAdminId, typeRequest, requestedChapterId)

	newRequest := models.Request{
		Id:                 id,
		Type:               typeRequest,
		RequestingAdmin:    requestingAdminId,
		RequestedChapterId: requestedChapterId,
	}

	_, err := storage.RegisterRequest(db, newRequest)

	if err != nil {
		log.Println("ошибка регестрирования запроса")
		return 0, err
	}

	admin, err := storage.SelectAdminWithId(db, requestingAdminId)

	if err != nil {
		log.Println("ошибка обноружения админа")
		return 0, err
	}

	admin.RequestSent = append(admin.RequestSent, id)

	_, err = storage.UpdateAdmin(db, admin.Id, admin)

	if err != nil {
		return 0, err
	}

	admins, err := storage.SelectAllSupeAdmins(db)

	for _, admin := range admins {
		ad, err := storage.SelectAdminWithId(db, admin.Id)

		if err == nil {
			ad.RequestsReceived = append(ad.RequestsReceived, id)
		}

		_, _ = storage.UpdateAdmin(db, admin.Id, ad)
	}

	log.Println(typeRequest, "typeRequest")

	if typeRequest == 1 && requestedChapterId != 0 {

		chapter, err := storage.SelectChapterWIthId(db, requestedChapterId)

		if err != nil {
			return 0, err
		}

		chapter.Status = 2

		_, err = storage.UpdateChapter(db, chapter.Id, chapter)

		if err != nil {
			return 0, err
		}

		log.Println("yовый статус", chapter.Status)
	}

	return id, nil
}
