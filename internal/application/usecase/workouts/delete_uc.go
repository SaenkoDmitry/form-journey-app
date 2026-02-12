package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type DeleteUseCase struct {
	workoutsRepo  workouts.Repo
	setsRepo      sets.Repo
	exercisesRepo exercises.Repo
}

func NewDeleteUseCase(workoutsRepo workouts.Repo, setsRepo sets.Repo, exercisesRepo exercises.Repo) *DeleteUseCase {
	return &DeleteUseCase{workoutsRepo: workoutsRepo, setsRepo: setsRepo, exercisesRepo: exercisesRepo}
}

func (uc *DeleteUseCase) Name() string {
	return "Удаление тренировки"
}

func (uc *DeleteUseCase) Execute(workoutID int64) error {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return err
	}

	for _, exercise := range workoutDay.Exercises {
		deleteErr := uc.setsRepo.DeleteAllBy(exercise.ID)
		if deleteErr != nil {
			return deleteErr
		}
	}

	err = uc.exercisesRepo.DeleteByWorkout(workoutID)
	if err != nil {
		return err
	}

	err = uc.workoutsRepo.Delete(&workoutDay)
	if err != nil {
		return err
	}
	return nil
}
