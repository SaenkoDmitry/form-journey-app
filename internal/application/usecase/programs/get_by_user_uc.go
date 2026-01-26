package programs

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetByUserUseCase struct {
	programsRepo programs.Repo
	usersRepo    users.Repo
}

func NewGetByUserUseCase(
	programsRepo programs.Repo,
	usersRepo users.Repo,
) *GetByUserUseCase {
	return &GetByUserUseCase{
		programsRepo: programsRepo,
		usersRepo:    usersRepo,
	}
}

func (uc *GetByUserUseCase) Name() string {
	return "Загрузить программу пользователя"
}

func (uc *GetByUserUseCase) Execute(chatID int64) (*models.WorkoutProgram, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	program, err := uc.programsRepo.Get(*user.ActiveProgramID)
	if err != nil {
		return nil, err
	}

	return &program, nil
}
