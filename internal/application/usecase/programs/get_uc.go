package programs

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type GetUseCase struct {
	usersRepo        users.Repo
	programsRepo     programs.Repo
	exerciseTypeRepo exercisetypes.Repo
}

func NewGetUseCase(
	programsRepo programs.Repo,
	exerciseTypeRepo exercisetypes.Repo,
	usersRepo users.Repo,
) *GetUseCase {
	return &GetUseCase{
		usersRepo:        usersRepo,
		programsRepo:     programsRepo,
		exerciseTypeRepo: exerciseTypeRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Редактировать программу"
}

func (uc *GetUseCase) ExecuteByChatID(programID, chatID int64) (*dto.ProgramDTO, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return uc.Execute(programID, user.ID)
}

func (uc *GetUseCase) Execute(programID, userID int64) (*dto.ProgramDTO, error) {
	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	programObj, err := uc.programsRepo.Get(programID)
	if err != nil {
		return nil, err
	}

	return mapProgramDTO(programObj, user), nil
}
