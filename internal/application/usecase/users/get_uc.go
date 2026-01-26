package users

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetUseCase struct {
	usersRepo users.Repo
}

func NewGetUseCase(usersRepo users.Repo) *GetUseCase {
	return &GetUseCase{
		usersRepo: usersRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Найти пользователя в системе"
}

func (uc *GetUseCase) Execute(chatID int64) (*models.User, error) {
	return uc.usersRepo.GetByChatID(chatID)
}
