package sets

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator"
	summarysvc "github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

type UpdateNextUseCase struct {
	usersRepo              users.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
	workoutsRepo           workouts.Repo
	exercisesRepo          exercises.Repo
	summaryService         summarysvc.Service
	docGeneratorService    docgenerator.Service
	setsRepo               sets.Repo
}

func NewUpdateNextUseCase(
	setsRepo sets.Repo,
	exercisesRepo exercises.Repo,
) *UpdateNextUseCase {
	return &UpdateNextUseCase{
		setsRepo:      setsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *UpdateNextUseCase) Name() string {
	return "Изменить следующий повтор"
}

var (
	NotFoundSetErr = errors.New("not found set")
)

func (uc *UpdateNextUseCase) Execute(exerciseID int64, newSetDTO *dto.NewSet) (int64, error) {
	exercise, _ := uc.exercisesRepo.Get(exerciseID)
	nextSet := exercise.NextSet()
	if nextSet.ID == 0 {
		return 0, NotFoundSetErr
	}

	if newSetDTO.NewReps > 0 {
		if int(newSetDTO.NewReps) != nextSet.Reps {
			nextSet.FactReps = int(newSetDTO.NewReps)
		} else {
			nextSet.FactReps = 0
		}
	}

	if newSetDTO.NewWeight > 0 {
		if float32(newSetDTO.NewWeight) != nextSet.Weight {
			nextSet.FactWeight = float32(newSetDTO.NewWeight)
		} else {
			nextSet.FactWeight = float32(0)
		}
	}

	if newSetDTO.NewMinutes > 0 {
		if int(newSetDTO.NewMinutes) != nextSet.Minutes {
			nextSet.FactMinutes = int(newSetDTO.NewMinutes)
		} else {
			nextSet.FactMinutes = 0
		}
	}

	if newSetDTO.NewMeters > 0 {
		if int(newSetDTO.NewMeters) != nextSet.Meters {
			nextSet.FactMeters = int(newSetDTO.NewMeters)
		} else {
			nextSet.FactMeters = 0
		}
	}

	err := uc.setsRepo.Save(&nextSet)
	if err != nil {
		return 0, err
	}

	return exercise.WorkoutDayID, nil
}
