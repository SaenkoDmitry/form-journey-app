package workouts

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type ShowProgressUseCase struct {
	workoutsRepo workouts.Repo
	sessionsRepo sessions.Repo
}

func NewShowProgressUseCase(
	workoutsRepo workouts.Repo,
	sessionsRepo sessions.Repo,
) *ShowProgressUseCase {
	return &ShowProgressUseCase{
		workoutsRepo: workoutsRepo,
		sessionsRepo: sessionsRepo,
	}
}

func (uc *ShowProgressUseCase) Name() string {
	return "Показать прогресс тренировки"
}

var (
	ErrWorkoutNotFound = errors.New("тренировка не найдена")
)

func (uc *ShowProgressUseCase) Execute(workoutID int64) (*dto.WorkoutProgress, error) {
	w, err := uc.workoutsRepo.Get(workoutID)
	if err != nil || w.ID == 0 {
		return nil, ErrWorkoutNotFound
	}

	totalExercises := len(w.Exercises)
	totalSets := 0
	completedExercises := 0
	completedSets := 0

	for _, exercise := range w.Exercises {
		setsCount := len(exercise.Sets)
		totalSets += setsCount

		done := exercise.CompletedSets()
		completedSets += done

		if done == setsCount && setsCount > 0 {
			completedExercises++
		}
	}

	progress := 0
	if totalSets > 0 {
		progress = (completedSets * 100) / totalSets
	}

	var remaining *int
	if w.EndedAt == nil && completedSets > 0 {
		elapsed := time.Since(w.StartedAt)
		setsPerMinute := float64(completedSets) / elapsed.Minutes()

		if setsPerMinute > 0 {
			left := totalSets - completedSets
			min := int(float64(left) / setsPerMinute)
			remaining = &min
		}
	}

	session, _ := uc.sessionsRepo.GetByWorkoutID(workoutID)

	return &dto.WorkoutProgress{
		Workout:            &w,
		TotalExercises:     totalExercises,
		CompletedExercises: completedExercises,
		TotalSets:          totalSets,
		CompletedSets:      completedSets,
		ProgressPercent:    progress,
		RemainingMin:       remaining,
		SessionStarted:     session.IsActive,
	}, nil
}
