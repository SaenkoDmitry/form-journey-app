package programs

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
)

type FindAllByUserUseCase struct {
	programsRepo programs.Repo
	usersRepo    users.Repo
}

func NewFindAllByUserUseCase(
	programsRepo programs.Repo,
	usersRepo users.Repo,
) *FindAllByUserUseCase {
	return &FindAllByUserUseCase{
		programsRepo: programsRepo,
		usersRepo:    usersRepo,
	}
}

func (uc *FindAllByUserUseCase) Name() string {
	return "Управление программами"
}

var (
	NoProgramsErr = errors.New("no training programs")
)

func (uc *FindAllByUserUseCase) Execute(chatID int64) (*dto.GetAllPrograms, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	programObjs, err := uc.programsRepo.FindAll(user.ID)
	if err != nil {
		return nil, err
	}

	if len(programObjs) == 0 {
		return nil, NoProgramsErr
	}
	return &dto.GetAllPrograms{
		User:     user,
		Programs: programObjs,
	}, nil
}
