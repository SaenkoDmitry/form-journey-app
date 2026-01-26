package workouts

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"time"
)

type StartUseCase struct {
	workoutsRepo  workouts.Repo
	exercisesRepo exercises.Repo
	usersRepo     users.Repo
	sessionsRepo  sessions.Repo
}

func NewStartUseCase(workoutsRepo workouts.Repo, sessionsRepo sessions.Repo) *StartUseCase {
	return &StartUseCase{workoutsRepo: workoutsRepo, sessionsRepo: sessionsRepo}
}

func (uc *StartUseCase) Name() string {
	return "Начать тренировку"
}

var (
	NotFoundSpecificErr = errors.New("workout not found")
	AlreadyCompletedErr = errors.New("already completed")
)

func (uc *StartUseCase) Execute(workoutID int64) (*dto.StartWorkout, error) {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	if workoutDay.ID == 0 {
		return nil, NotFoundSpecificErr
	}

	if workoutDay.Completed {
		return nil, AlreadyCompletedErr
	}

	session := models.WorkoutSession{
		WorkoutDayID:         workoutDay.ID,
		StartedAt:            time.Now(),
		IsActive:             true,
		CurrentExerciseIndex: 0,
	}
	if err = uc.sessionsRepo.Create(&session); err != nil {
		return nil, err
	}

	return &dto.StartWorkout{}, nil
}
