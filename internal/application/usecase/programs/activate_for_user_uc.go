package programs

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type ActivateUseCase struct {
	usersRepo users.Repo
}

func NewActivateUseCase(
	usersRepo users.Repo,
) *ActivateUseCase {
	return &ActivateUseCase{
		usersRepo: usersRepo,
	}
}

func (uc *ActivateUseCase) Name() string {
	return "Активировать программу"
}

func (uc *ActivateUseCase) Execute(userID, programID int64) error {
	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return err
	}

	*user.ActiveProgramID = programID
	err = uc.usersRepo.Save(user)
	if err != nil {
		return err
	}

	return nil
}
