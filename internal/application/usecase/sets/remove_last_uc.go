package sets

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
)

type RemoveLastUseCase struct {
	setsRepo      sets.Repo
	exercisesRepo exercises.Repo
}

func NewRemoveLastUseCase(
	setsRepo sets.Repo,
	exercisesRepo exercises.Repo,
) *RemoveLastUseCase {
	return &RemoveLastUseCase{
		setsRepo:      setsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *RemoveLastUseCase) Name() string {
	return "Удалить подход"
}

var (
	AddOneMoreExerciseToDeleteErr = errors.New(messages.AddOneMoreExerciseToDelete)
	YouCannotDeleteOneExerciseErr = errors.New(messages.YouCannotDeleteOneOfSet)
)

func (uc *RemoveLastUseCase) Execute(exerciseID int64) (*dto.RemoveLastSet, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil || len(exercise.Sets) == 0 {
		return nil, AddOneMoreExerciseToDeleteErr
	}
	if len(exercise.Sets) == 1 {
		return nil, YouCannotDeleteOneExerciseErr
	}

	lastSet := exercise.Sets[len(exercise.Sets)-1]
	err = uc.setsRepo.Delete(lastSet.ID)
	if err != nil {
		fmt.Println("cannot remove set:", err.Error())
		return nil, err
	}

	return &dto.RemoveLastSet{WorkoutID: exercise.WorkoutDayID}, nil
}
