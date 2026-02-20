package workouts

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
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

func (uc *CreateUseCase) Execute(userID, dayTypeID int64) (*dto.CreateWorkout, error) {
	workout := &models.WorkoutDay{
		UserID:           userID,
		WorkoutDayTypeID: dayTypeID,
		StartedAt:        time.Now(),
		Completed:        false,
	}
	if err := uc.workoutsRepo.Create(workout); err != nil {
		return nil, err
	}

	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	exerciseObjs, err := uc.buildExercises(workout.ID, dayTypeID, user)
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

func (uc *CreateUseCase) buildExercises(workoutID int64, dayTypeID int64, user *models.User) ([]models.Exercise, error) {
	activeProgramID := *user.ActiveProgramID
	return uc.createExercisesFromPresets(workoutID, dayTypeID, activeProgramID)
}

func (uc *CreateUseCase) createExercisesFromPresets(workoutDayID, dayTypeID, activeProgramID int64) ([]models.Exercise, error) {
	method := "createExercisesFromPresets"
	fmt.Printf("%s: берем настройки количества повторений и веса из preset-ов\n", method)

	objs := make([]models.Exercise, 0)
	dayType, err := uc.dayTypesRepo.Get(dayTypeID)
	if err != nil {
		return nil, err
	}

	for index, presetEx := range utils.SplitPreset(dayType.Preset) {

		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: presetEx.ID,
			Index:          index,
		}

		if prevEx, prevErr := uc.exercisesRepo.FindPreviousByType(presetEx.ID, activeProgramID); prevErr == nil {
			newExercise.Sets = prevEx.CloneSets()
		} else {
			for idx2, set := range presetEx.Sets {
				newSet := models.Set{Index: idx2}
				if set.Minutes > 0 {
					newSet.Minutes = set.Minutes
				} else {
					newSet.Reps = set.Reps
					newSet.Weight = set.Weight
				}
				newExercise.Sets = append(newExercise.Sets, newSet)
			}
		}

		objs = append(objs, newExercise)
	}
	return objs, nil
}

func (uc *CreateUseCase) createExercisesFromLastWorkout(workoutDayID, previousWorkoutID, activeProgramID int64) ([]models.Exercise, error) {
	method := "createExercisesFromLastWorkout"
	fmt.Printf("%s: берем настройки количества повторений и веса из последней тренировки: %d\n", method, previousWorkoutID)

	previousExercises, err := uc.exercisesRepo.FindAllByWorkoutID(previousWorkoutID)
	if err != nil {
		return nil, err
	}
	objs := make([]models.Exercise, 0)
	for _, prevWorkoutEx := range previousExercises {
		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: prevWorkoutEx.ExerciseTypeID,
			Index:          prevWorkoutEx.Index,
		}

		if prevEx, prevErr := uc.exercisesRepo.FindPreviousByType(prevWorkoutEx.ExerciseTypeID, activeProgramID); prevErr == nil {
			newExercise.Sets = buildSetsFrom(prevEx)
		} else {
			newExercise.Sets = buildSetsFrom(prevWorkoutEx)
		}
		objs = append(objs, newExercise)
	}
	return objs, nil
}

func buildSetsFrom(previousEx models.Exercise) []models.Set {
	sets := make([]models.Set, 0, len(previousEx.Sets))
	for _, set := range previousEx.Sets {
		newSet := models.Set{
			Reps:    set.GetRealReps(),
			Weight:  set.GetRealWeight(),
			Minutes: set.GetRealMinutes(),
			Meters:  set.GetRealMeters(),
			Index:   set.Index,
		}
		switch {
		case newSet.Reps == 0 && previousEx.ExerciseType.ContainsReps():
			newSet.Reps = constants.DefaultReps
		case newSet.Weight == 0 && previousEx.ExerciseType.ContainsWeight():
			newSet.Weight = constants.DefaultWeight
		case newSet.Minutes == 0 && previousEx.ExerciseType.ContainsMinutes():
			newSet.Minutes = constants.DefaultMinutes
		case newSet.Meters == 0 && previousEx.ExerciseType.ContainsMeters():
			newSet.Meters = constants.DefaultMeters
		}
		sets = append(sets, newSet)
	}
	return sets
}
