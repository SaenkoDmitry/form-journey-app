package session

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"time"
)

type MoveToUseCase struct {
	sessionsRepo  sessions.Repo
	exercisesRepo exercises.Repo
}

func NewMoveToExerciseUseCase(
	sessionsRepo sessions.Repo,
	exercisesRepo exercises.Repo,
) *MoveToUseCase {
	return &MoveToUseCase{
		sessionsRepo:  sessionsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *MoveToUseCase) Name() string {
	return "Перейти к другому упражнению"
}

var (
	NoExercisesInWorkout        = errors.New("no exercises in workout")
	NoEarlierExercisesInWorkout = errors.New("no earlier exercises in workout")
	YouCompletedAllExercises    = errors.New("you completed all exercises in workout")
)

func (uc *MoveToUseCase) Execute(workoutID int64, next bool) error {
	session, err := uc.sessionsRepo.GetByWorkoutID(workoutID)
	if err != nil {
		// или ошибка?
		session = models.WorkoutSession{
			WorkoutDayID:         workoutID,
			StartedAt:            time.Now(),
			IsActive:             true,
			CurrentExerciseIndex: 0,
		}
		if err = uc.sessionsRepo.Create(&session); err != nil {
			return err
		}
	}

	exerciseObjs, err := uc.exercisesRepo.FindAllByWorkoutID(workoutID)
	if err != nil {
		return err
	}

	if len(exerciseObjs) == 0 {
		return NoExercisesInWorkout
	}

	if next {
		session.CurrentExerciseIndex++
	} else {
		session.CurrentExerciseIndex--
	}

	if session.CurrentExerciseIndex < 0 {
		session.CurrentExerciseIndex = 0
		return NoEarlierExercisesInWorkout
	}

	if session.CurrentExerciseIndex >= len(exerciseObjs) {
		session.CurrentExerciseIndex = 0
		return YouCompletedAllExercises
	}

	err = uc.sessionsRepo.Save(&session)
	if err != nil {
		return err
	}

	return nil
}
