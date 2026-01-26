package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"time"
)

type FinishUseCase struct {
	workoutsRepo workouts.Repo
	dayTypesRepo daytypes.Repo
	sessionsRepo sessions.Repo
}

func NewFinishUseCase(workoutsRepo workouts.Repo, sessionsRepo sessions.Repo) *FinishUseCase {
	return &FinishUseCase{workoutsRepo: workoutsRepo, sessionsRepo: sessionsRepo}
}

func (uc *FinishUseCase) Name() string {
	return "Завершение тренировки"
}

func (uc *FinishUseCase) Execute(workoutID int64) (*dto.FinishWorkout, error) {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	workoutDay.Completed = true
	workoutDay.EndedAt = &now

	if err = uc.workoutsRepo.Save(&workoutDay); err != nil {
		return nil, err
	}

	err = uc.sessionsRepo.UpdateIsActive(workoutID, false)
	if err != nil {
		return nil, err
	}

	return &dto.FinishWorkout{WorkoutID: workoutDay.ID}, nil
}
