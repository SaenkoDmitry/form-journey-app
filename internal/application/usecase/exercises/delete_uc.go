package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
)

type DeleteUseCase struct {
	exercisesRepo exercises.Repo
}

func NewDeleteUseCase(exercisesRepo exercises.Repo) *DeleteUseCase {
	return &DeleteUseCase{
		exercisesRepo: exercisesRepo,
	}
}

func (uc *DeleteUseCase) Name() string {
	return "Удалить упражнение"
}

func (uc *DeleteUseCase) Execute(exerciseID int64) (int64, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return 0, err
	}
	workoutID := exercise.WorkoutDayID

	err = uc.exercisesRepo.Delete(exerciseID)
	if err != nil {
		return 0, err
	}
	return workoutID, nil
}
