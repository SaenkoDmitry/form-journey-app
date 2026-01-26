package programs

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
)

type RenameUseCase struct {
	programsRepo programs.Repo
}

func NewRenameUseCase(
	programsRepo programs.Repo,
) *RenameUseCase {
	return &RenameUseCase{
		programsRepo: programsRepo,
	}
}

func (uc *RenameUseCase) Name() string {
	return "Переименовать программу"
}

func (uc *RenameUseCase) Execute(programID int64, newName string) error {
	program, err := uc.programsRepo.Get(programID)
	if err != nil {
		return err
	}
	program.Name = newName
	if err = uc.programsRepo.Save(&program); err != nil {
		return err
	}
	return nil
}
