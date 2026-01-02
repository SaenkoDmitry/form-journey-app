package users

import (
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	GetUser(chatID int64, username string) (*models.User, error)
	GetUserByChatID(chatID int64) *models.User
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) GetUser(chatID int64, username string) (*models.User, error) {
	var user models.User
	result := u.db.Where("chat_id = ?", chatID).First(&user)

	if result.Error != nil {
		user = models.User{
			ChatID:    chatID,
			Username:  username,
			CreatedAt: time.Now(),
		}
		u.db.Create(&user)
	}

	return &user, nil
}

func (u *repoImpl) GetUserByChatID(chatID int64) *models.User {
	var user models.User
	u.db.Where("chat_id = ?", chatID).First(&user)
	return &user
}
