package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
)

type ConfirmDeleteUseCase struct {
	exercisesRepo     exercises.Repo
	exerciseTypesRepo exercisetypes.Repo
}

func NewConfirmDeleteUseCase(exerciseTypesRepo exercisetypes.Repo, exercisesRepo exercises.Repo) *ConfirmDeleteUseCase {
	return &ConfirmDeleteUseCase{
		exercisesRepo:     exercisesRepo,
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *ConfirmDeleteUseCase) Name() string {
	return "Удалить упражнение"
}

func (uc *ConfirmDeleteUseCase) Execute(exerciseID int64) (*dto.ConfirmDeleteExercise, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return nil, err
	}

	exerciseObj, err := uc.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.ConfirmDeleteExercise{
		Exercise:    exercise,
		ExerciseObj: exerciseObj,
	}, nil
}
