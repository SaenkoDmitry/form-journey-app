package pushsubscriptions

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/pushsubscriptions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type CreateUseCase struct {
	pushSubscriptionsRepo pushsubscriptions.Repo
	usersRepo             users.Repo
}

func NewCreateUseCase(
	pushSubscriptionsRepo pushsubscriptions.Repo,
	usersRepo users.Repo,
) *CreateUseCase {
	return &CreateUseCase{
		pushSubscriptionsRepo: pushSubscriptionsRepo,
		usersRepo:             usersRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Создать подписку"
}

func (uc *CreateUseCase) Execute(userID int64, sub dto.PushSubscription) error {
	return uc.pushSubscriptionsRepo.Create(userID, sub)
}
