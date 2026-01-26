package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
)

type GetUseCase struct {
	exercisesRepo exercises.Repo
}

func NewGetUseCase(
	exercisesRepo exercises.Repo,
) *GetUseCase {
	return &GetUseCase{
		exercisesRepo: exercisesRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Показать данные об упражнении"
}

func (uc *GetUseCase) Execute(exerciseID int64) (*dto.GetExercise, error) {
	ex, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return nil, err
	}

	return &dto.GetExercise{Exercise: ex}, nil
}
