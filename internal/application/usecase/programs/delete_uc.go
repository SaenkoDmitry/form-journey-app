package programs

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type DeleteUseCase struct {
	programsRepo programs.Repo
	usersRepo    users.Repo
}

func NewDeleteUseCase(
	programsRepo programs.Repo,
	usersRepo users.Repo,
) *DeleteUseCase {
	return &DeleteUseCase{
		programsRepo: programsRepo,
		usersRepo:    usersRepo,
	}
}

func (uc *DeleteUseCase) Name() string {
	return "Удалить программу"
}

var (
	CannotDeleteCurrentProgramErr = errors.New("cannot delete current program")
)

func (uc *DeleteUseCase) Execute(chatID, programID int64) error {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return err
	}

	if *user.ActiveProgramID == programID {
		return CannotDeleteCurrentProgramErr
	}

	program, err := uc.programsRepo.Get(programID)
	if err != nil {
		return err
	}

	err = uc.programsRepo.Delete(&program)
	if err != nil {
		return err
	}

	return nil
}
