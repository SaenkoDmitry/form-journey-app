package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type ConfirmFinishUseCase struct {
	workoutsRepo workouts.Repo
	dayTypesRepo daytypes.Repo
}

func NewConfirmFinishUseCase(workoutsRepo workouts.Repo, dayTypesRepo daytypes.Repo) *ConfirmFinishUseCase {
	return &ConfirmFinishUseCase{workoutsRepo: workoutsRepo, dayTypesRepo: dayTypesRepo}
}

func (uc *ConfirmFinishUseCase) Name() string {
	return "Завершение тренировки"
}

func (uc *ConfirmFinishUseCase) Execute(workoutID int64) (*dto.ConfirmFinishWorkout, error) {
	workoutDay, err := uc.workoutsRepo.Get(workoutID)
	if err != nil {
		return nil, err
	}

	dayType, err := uc.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.ConfirmFinishWorkout{
		DayType: dayType,
	}, nil
}
