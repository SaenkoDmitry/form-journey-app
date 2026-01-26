package sets

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"time"
)

type CompleteUseCase struct {
	setsRepo          sets.Repo
	exercisesRepo     exercises.Repo
	exerciseTypesRepo exercisetypes.Repo
}

func NewCompleteUseCase(
	setsRepo sets.Repo,
	exercisesRepo exercises.Repo,
	exerciseTypesRepo exercisetypes.Repo,
) *CompleteUseCase {
	return &CompleteUseCase{
		setsRepo:          setsRepo,
		exercisesRepo:     exercisesRepo,
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *CompleteUseCase) Name() string {
	return "Завершить подход"
}

var (
	DoNothingErr = errors.New("do nothing")
)

func (uc *CompleteUseCase) Execute(exerciseID int64) (*dto.CompleteSet, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return nil, err
	}

	nextSet := exercise.NextSet()

	if nextSet.ID != 0 {
		nextSet.Completed = true
		now := time.Now()
		nextSet.CompletedAt = &now
		uc.setsRepo.Save(&nextSet)
	} else {
		return nil, DoNothingErr
	}

	exerciseType, err := uc.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		return nil, DoNothingErr
	}

	return &dto.CompleteSet{
		NeedMoveToNext:  nextSet.ID == exercise.LastSet().ID,
		NeedShowCurrent: nextSet.ID != exercise.LastSet().ID,
		NeedStartTimer:  exerciseType.RestInSeconds > 0,
		WorkoutID:       exercise.WorkoutDayID,
		Seconds:         exerciseType.RestInSeconds,
	}, nil
}
