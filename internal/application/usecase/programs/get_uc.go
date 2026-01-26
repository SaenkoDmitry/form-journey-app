package programs

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
)

type GetUseCase struct {
	programsRepo     programs.Repo
	exerciseTypeRepo exercisetypes.Repo
}

func NewGetUseCase(
	programsRepo programs.Repo,
	exerciseTypeRepo exercisetypes.Repo,
) *GetUseCase {
	return &GetUseCase{
		programsRepo:     programsRepo,
		exerciseTypeRepo: exerciseTypeRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Редактировать программу"
}

func (uc *GetUseCase) Execute(programID int64) (*dto.GetProgram, error) {
	program, err := uc.programsRepo.Get(programID)
	if err != nil {
		return nil, err
	}

	exerciseTypesMap := make(map[int64]models.ExerciseType)
	exTypes, err := uc.exerciseTypeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, ex := range exTypes {
		exerciseTypesMap[ex.ID] = ex
	}

	return &dto.GetProgram{
		Program:          program,
		ExerciseTypesMap: exerciseTypesMap,
	}, nil
}
