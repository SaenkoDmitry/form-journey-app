package users

import (
	"errors"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Save(user *models.User) error
	Create(chatID int64, from *tgbotapi.User) (*models.User, error)
	GetByChatID(chatID int64) (*models.User, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Save(user *models.User) error {
	u.db.Save(user)
	return nil
}

func (u *repoImpl) Create(chatID int64, from *tgbotapi.User) (*models.User, error) {
	user := models.User{
		ChatID:       chatID,
		Username:     from.UserName,
		FirstName:    from.FirstName,
		LastName:     from.LastName,
		LanguageCode: from.LanguageCode,
		CreatedAt:    time.Now(),
	}
	err := u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})
	return &user, err
}

var (
	NotFoundUserErr = errors.New("not found user")
)

func (u *repoImpl) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User

	result := u.db.Where("chat_id = ?", chatID).First(&user)

	if result.Error != nil {
		return nil, NotFoundUserErr
	}
	return &user, nil
}
