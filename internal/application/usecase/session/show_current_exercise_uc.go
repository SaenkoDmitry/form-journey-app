package session

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
)

type ShowCurrentExerciseSessionUseCase struct {
	workoutsRepo      workouts.Repo
	sessionsRepo      sessions.Repo
	exerciseTypesRepo exercisetypes.Repo
	dayTypesRepo      daytypes.Repo
}

func NewShowCurrentExerciseUseCase(
	workoutsRepo workouts.Repo,
	sessionsRepo sessions.Repo,
	exerciseTypesRepo exercisetypes.Repo,
	dayTypesRepo daytypes.Repo,
) *ShowCurrentExerciseSessionUseCase {
	return &ShowCurrentExerciseSessionUseCase{
		workoutsRepo:      workoutsRepo,
		sessionsRepo:      sessionsRepo,
		exerciseTypesRepo: exerciseTypesRepo,
		dayTypesRepo:      dayTypesRepo,
	}
}

func (uc *ShowCurrentExerciseSessionUseCase) Name() string {
	return "Показать текущее упражнение"
}

var (
	NoExercisesErr      = errors.New("no exercises")
	NotFoundExerciseErr = errors.New("exercise not found")
)

func (uc *ShowCurrentExerciseSessionUseCase) Execute(workoutID int64) (*dto.CurrentExerciseSession, error) {
	workoutDay, _ := uc.workoutsRepo.Get(workoutID)
	if len(workoutDay.Exercises) == 0 {
		return nil, NoExercisesErr
	}

	session, err := uc.sessionsRepo.GetByWorkoutID(workoutID)
	if err != nil {
		return nil, err
	}

	exerciseIndex := session.CurrentExerciseIndex
	if exerciseIndex >= len(workoutDay.Exercises) {
		exerciseIndex = 0
	}

	exercise := workoutDay.Exercises[exerciseIndex]

	exerciseObj, err := uc.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		return nil, NotFoundExerciseErr
	}

	dayType, err := uc.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return nil, err
	}

	return &dto.CurrentExerciseSession{
		ExerciseIndex: exerciseIndex,
		WorkoutDay:    workoutDay,
		Exercise:      exercise,
		ExerciseObj:   exerciseObj,
		DayType:       dayType,
	}, nil
}
