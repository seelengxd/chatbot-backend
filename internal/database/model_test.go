package database

import (
	"backend/internal/configs"
	"backend/internal/model"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestCreateChat(t *testing.T) {
	SetupTestDb(t)
	t.Run("can create chat without booking id", func(t *testing.T) {
		db := GetDb()
		defer CleanUpTestDb(t, db)
		chat := model.Chat{TelegramChatId: 1}

		err := (&chat).Create(db)

		AssertNoErr(t, err)
	})
	t.Run("can create chat with booking id", func(t *testing.T) {
		db := GetDb()
		defer CleanUpTestDb(t, db)
		booking := model.Booking{LastName: "john", RoomNumber: "02-01"}
		chat := model.Chat{TelegramChatId: 1, Booking: &booking}

		AssertNoErr(t, chat.Create(db))
	})
}

func TestCreateMessage(t *testing.T) {
	SetupTestDb(t)
	db := GetDb()
	defer CleanUpTestDb(t, db)

	chat := model.Chat{TelegramChatId: 1}
	AssertNoErr(t, (&chat).Create(db))

	query := model.RequestQuery{
		Status:    model.StatusOngoing,
		Type:      model.TypeUnknown,
		BookingId: nil,
		Booking:   nil,
	}
	err := db.Model(&chat).Association("RequestQueries").Append(&query)
	AssertNoErr(t, err)

	user := model.User{
		Username:          "username",
		EncryptedPassword: "password",
	}
	AssertNoErr(t, user.Create(db))

	message := model.Message{
		TelegramMessageId: 1,
		By:                model.ByStaff,
		MessageBody:       "very cool message",
		Timestamp:         time.Now(),
		HotelStaffId:      user.ID,
		RequestQueryId:    query.ID,
	}
	AssertNoErr(t, message.Create(db))
}

func SetupTestDb(t testing.TB) {
	t.Helper()
	cfg, err := configs.GetConfig()
	if err != nil {
		panic(err)
	}
	SetupDb(cfg.GetTestDatabaseConfig())
}

func AssertNoErr(t testing.TB, err error) {
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
}

func CleanUpTestDb(t testing.TB, db *gorm.DB) {
	t.Helper()
	db.Unscoped().Where("1 = 1").Delete(&model.Message{})
	db.Unscoped().Where("1 = 1").Delete(&model.RequestQuery{})
	db.Unscoped().Where("1 = 1").Delete(&model.Booking{})
	db.Unscoped().Where("1 = 1").Delete(&model.Chat{})
	db.Unscoped().Where("1 = 1").Delete(&model.User{})

	fmt.Println("did u even run")
	fmt.Println(db.Error)
}