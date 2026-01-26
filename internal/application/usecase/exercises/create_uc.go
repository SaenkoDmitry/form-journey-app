package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"strings"
)

type CreateUseCase struct {
	exercisesRepo     exercises.Repo
	workoutsRepo      workouts.Repo
	exerciseTypesRepo exercisetypes.Repo
}

func NewCreateUseCase(exercisesRepo exercises.Repo, workoutsRepo workouts.Repo, exerciseTypesRepo exercisetypes.Repo) *CreateUseCase {
	return &CreateUseCase{
		exercisesRepo:     exercisesRepo,
		workoutsRepo:      workoutsRepo,
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Добавить в тренировку упражнение"
}

func (uc *CreateUseCase) Execute(workoutID int64, exerciseTypeID int64) (*dto.CreateExercise, error) {
	exerciseObj, err := uc.exerciseTypesRepo.Get(exerciseTypeID)
	if err != nil {
		return nil, err
	}

	workout, _ := uc.workoutsRepo.Get(workoutID)
	idx := 0
	if len(workout.Exercises) > 0 {
		lastExercise := workout.Exercises[len(workout.Exercises)-1]
		idx = lastExercise.Index + 1
	}

	newExercise := models.Exercise{
		ExerciseTypeID: exerciseObj.ID,
		Index:          idx,
		WorkoutDayID:   workoutID,
		Sets: []models.Set{
			{Index: 1}, // по дефолту один подход
		},
	}
	switch {
	case strings.Contains(exerciseObj.Units, "meters"):
		newExercise.Sets[0].Meters = 100
	case strings.Contains(exerciseObj.Units, "minutes"):
		newExercise.Sets[0].Minutes = 1
	case strings.Contains(exerciseObj.Units, "reps"):
		newExercise.Sets[0].Reps = 10
	case strings.Contains(exerciseObj.Units, "weight"):
		newExercise.Sets[0].Weight = 10
	}
	err = uc.exercisesRepo.Save(&newExercise)
	if err != nil {
		return nil, err
	}

	return &dto.CreateExercise{
		ExerciseObj: exerciseObj,
	}, nil
}
