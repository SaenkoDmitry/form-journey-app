package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
)

type ExerciseTypeListUseCase struct {
	exerciseTypesRepo exercisetypes.Repo
}

func NewExerciseTypeListUseCase(exerciseTypesRepo exercisetypes.Repo) *ExerciseTypeListUseCase {
	return &ExerciseTypeListUseCase{
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *ExerciseTypeListUseCase) Name() string {
	return "Список упражнений"
}

func (uc *ExerciseTypeListUseCase) Execute() (*dto.ExerciseTypeList, error) {
	exerciseTypes, err := uc.exerciseTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return &dto.ExerciseTypeList{
		ExerciseTypes: exerciseTypes,
	}, nil
}
