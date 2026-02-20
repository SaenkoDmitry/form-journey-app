package users

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetByIDUseCase struct {
	usersRepo users.Repo
}

func NewGetByIDUseCase(usersRepo users.Repo) *GetByIDUseCase {
	return &GetByIDUseCase{
		usersRepo: usersRepo,
	}
}

func (uc *GetByIDUseCase) Name() string {
	return "Найти пользователя в системе"
}

func (uc *GetByIDUseCase) Execute(userID int64) (*models.User, error) {
	return uc.usersRepo.GetByID(userID)
}
