package programs

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"strings"
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
	CannotDeleteCurrentProgramErr  = errors.New("Не могу удалить активную программу")
	CannotDeleteAlreadyUsedProgram = errors.New("Не могу удалить программу, которая уже есть в истории тренировок")
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
		if strings.Contains(err.Error(), "update or delete on table \"workout_day_types\" violates foreign key constraint") {
			return CannotDeleteAlreadyUsedProgram
		}
		return err
	}

	return nil
}
