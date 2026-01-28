package sets

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/docgenerator"
	summarysvc "github.com/SaenkoDmitry/training-tg-bot/internal/service/summary"
)

type AddOneMoreUseCase struct {
	usersRepo              users.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
	workoutsRepo           workouts.Repo
	exercisesRepo          exercises.Repo
	summaryService         summarysvc.Service
	docGeneratorService    docgenerator.Service
	setsRepo               sets.Repo
}

func NewAddOneMoreUseCase(
	setsRepo sets.Repo,
	exercisesRepo exercises.Repo,
) *AddOneMoreUseCase {
	return &AddOneMoreUseCase{
		setsRepo:      setsRepo,
		exercisesRepo: exercisesRepo,
	}
}

func (uc *AddOneMoreUseCase) Name() string {
	return "Добавить подход"
}

func (uc *AddOneMoreUseCase) Execute(exerciseID int64) (*dto.AddOneMoreSet, error) {
	exercise, err := uc.exercisesRepo.Get(exerciseID)
	if err != nil {
		return nil, err
	}

	var nextSet *models.Set
	if len(exercise.Sets) > 0 {
		lastSet := exercise.Sets[len(exercise.Sets)-1]
		nextSet = &models.Set{
			ExerciseID: exercise.ID,
			Reps:       lastSet.GetRealReps(),
			Weight:     lastSet.GetRealWeight(),
			Minutes:    lastSet.GetRealMinutes(),
			Meters:     lastSet.GetRealMeters(),
			Index:      lastSet.Index + 1,
		}
	} else {
		nextSet = &models.Set{
			ExerciseID: exercise.ID,
		}
	}

	err = uc.setsRepo.Save(nextSet)
	if err != nil {
		fmt.Println("cannot create set:", err.Error())
		return nil, err
	}

	return &dto.AddOneMoreSet{
		WorkoutID: exercise.WorkoutDayID,
	}, nil
}
