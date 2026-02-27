package session

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"time"
)

type MoveToCertainUseCase struct {
	sessionsRepo  sessions.Repo
	exercisesRepo exercises.Repo
}

func NewMoveToCertainUseCase(
	sessionsRepo sessions.Repo,
	exercisesRepo exercises.Repo,
) *MoveToCertainUseCase {
	return &MoveToCertainUseCase{
		sessionsRepo:  sessionsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *MoveToCertainUseCase) Name() string {
	return "Перейти к упражнению"
}

func (uc *MoveToCertainUseCase) Execute(workoutID int64, exerciseIndex int) error {
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

	if exerciseIndex < 0 {
		session.CurrentExerciseIndex = 0
	}

	if exerciseIndex >= len(exerciseObjs) {
		session.CurrentExerciseIndex = len(exerciseObjs) - 1
	}

	session.CurrentExerciseIndex = exerciseIndex

	err = uc.sessionsRepo.Save(&session)
	if err != nil {
		return err
	}

	return nil
}
