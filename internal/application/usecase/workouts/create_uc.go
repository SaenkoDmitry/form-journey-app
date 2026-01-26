package workouts

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	"time"
)

type CreateUseCase struct {
	workoutsRepo  workouts.Repo
	exercisesRepo exercises.Repo
	usersRepo     users.Repo
	dayTypesRepo  daytypes.Repo
}

func NewCreateUseCase(workoutsRepo workouts.Repo, exercisesRepo exercises.Repo, usersRepo users.Repo, dayTypesRepo daytypes.Repo) *CreateUseCase {
	return &CreateUseCase{workoutsRepo: workoutsRepo, exercisesRepo: exercisesRepo, usersRepo: usersRepo, dayTypesRepo: dayTypesRepo}
}

func (uc *CreateUseCase) Name() string {
	return "Создание тренировки"
}

func (uc *CreateUseCase) Execute(chatID, dayTypeID int64) (*dto.CreateWorkout, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	workout := &models.WorkoutDay{
		UserID:           user.ID,
		WorkoutDayTypeID: dayTypeID,
		StartedAt:        time.Now(),
		Completed:        false,
	}
	if err = uc.workoutsRepo.Create(workout); err != nil {
		return nil, err
	}

	exerciseObjs, err := uc.buildExercises(workout.ID, dayTypeID, user.ID)
	if err != nil {
		return nil, err
	}

	err = uc.exercisesRepo.CreateBatch(exerciseObjs)
	if err != nil {
		return nil, err
	}

	return &dto.CreateWorkout{
		WorkoutID: workout.ID,
	}, nil
}

func (uc *CreateUseCase) buildExercises(workoutID int64, dayTypeID int64, userID int64) ([]models.Exercise, error) {
	previousWorkout, err := uc.workoutsRepo.FindPreviousByType(userID, dayTypeID)
	if err != nil {
		return uc.createExercisesFromPresets(workoutID, dayTypeID)
	}
	return uc.createExercisesFromLastWorkout(workoutID, previousWorkout.ID)
}

func (uc *CreateUseCase) createExercisesFromPresets(workoutDayID, dayTypeID int64) ([]models.Exercise, error) {
	method := "createExercisesFromPresets"
	fmt.Printf("%s: берем настройки количества повторений и веса из preset-ов\n", method)

	objs := make([]models.Exercise, 0)
	dayType, err := uc.dayTypesRepo.Get(dayTypeID)
	if err != nil {
		return nil, err
	}

	for index, exerciseType := range utils.SplitPreset(dayType.Preset) {
		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: exerciseType.ID,
			Index:          index,
		}
		for idx2, set := range exerciseType.Sets {
			newSet := models.Set{Index: idx2}
			if set.Minutes > 0 {
				newSet.Minutes = set.Minutes
			} else {
				newSet.Reps = set.Reps
				newSet.Weight = set.Weight
			}
			newExercise.Sets = append(newExercise.Sets, newSet)
		}
		objs = append(objs, newExercise)
	}
	return objs, nil
}

func (uc *CreateUseCase) createExercisesFromLastWorkout(workoutDayID, previousWorkoutID int64) ([]models.Exercise, error) {
	method := "createExercisesFromLastWorkout"
	fmt.Printf("%s: берем настройки количества повторений и веса из последней тренировки: %d\n", method, previousWorkoutID)

	previousExercises, err := uc.exercisesRepo.FindAllByWorkoutID(previousWorkoutID)
	if err != nil {
		return nil, err
	}
	objs := make([]models.Exercise, 0)
	for _, exercise := range previousExercises {
		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: exercise.ExerciseTypeID,
			Index:          exercise.Index,
		}
		for _, set := range exercise.Sets {
			newSet := models.Set{
				Reps:    set.GetRealReps(),
				Weight:  set.GetRealWeight(),
				Minutes: set.GetRealMinutes(),
				Meters:  set.GetRealMeters(),
				Index:   set.Index,
			}
			newExercise.Sets = append(newExercise.Sets, newSet)
		}
		objs = append(objs, newExercise)
	}
	return objs, nil
}
