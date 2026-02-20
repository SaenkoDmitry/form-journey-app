package pushsubscriptions

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/pushsubscriptions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type DeleteUseCase struct {
	pushSubscriptionsRepo pushsubscriptions.Repo
	usersRepo             users.Repo
}

func NewDeleteUseCase(
	pushSubscriptionsRepo pushsubscriptions.Repo,
	usersRepo users.Repo,
) *DeleteUseCase {
	return &DeleteUseCase{
		pushSubscriptionsRepo: pushSubscriptionsRepo,
		usersRepo:             usersRepo,
	}
}

func (uc *DeleteUseCase) Name() string {
	return "Удалить подписку"
}

func (uc *DeleteUseCase) Execute(userID int64, endpoint string) error {
	return uc.pushSubscriptionsRepo.Delete(userID, endpoint)
}
