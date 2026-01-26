package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type ConfirmDeleteUseCase struct {
	workoutsRepo workouts.Repo
	dayTypesRepo daytypes.Repo
}

func NewConfirmDeleteUseCase(workoutsRepo workouts.Repo, dayTypesRepo daytypes.Repo) *ConfirmDeleteUseCase {
	return &ConfirmDeleteUseCase{workoutsRepo: workoutsRepo, dayTypesRepo: dayTypesRepo}
}

func (uc *ConfirmDeleteUseCase) Name() string {
	return "Удаление тренировки"
}

func (uc *ConfirmDeleteUseCase) Execute(workoutID int64) (*dto.ConfirmDeleteWorkout, error) {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	dayType, err := uc.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.ConfirmDeleteWorkout{
		WorkoutID:   workoutDay.ID,
		DayTypeName: dayType.Name,
	}, nil
}
