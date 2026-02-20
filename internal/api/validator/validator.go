package validator

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/errorslist"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
)

func checkAccess(userID, entityID int64) error {
	if userID != entityID {
		return errorslist.ErrAccessDenied
	}
	return nil
}

func ValidateAccessToProgram(container *usecase.Container, userID int64, programID int64) error {
	program, err := container.GetProgramUC.Execute(programID, userID)
	if err != nil {
		return errorslist.ErrInternalMsg
	}
	return checkAccess(program.UserID, userID)
}

func ValidateAccessToExercise(container *usecase.Container, userID int64, exerciseID int64) error {
	ex, err := container.GetExerciseUC.Execute(exerciseID)
	if err != nil {
		return errorslist.ErrInternalMsg
	}
	return checkAccess(ex.Exercise.WorkoutDay.UserID, userID)
}

func ValidateAccessToWorkout(container *usecase.Container, userID int64, workoutID int64) error {
	progress, err := container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		return errorslist.ErrInternalMsg
	}
	return checkAccess(progress.Workout.UserID, userID)
}

func ValidateAccessToMeasurement(container *usecase.Container, userID int64, measurementID int64) error {
	measurement, err := container.GetMeasurementByIDUC.Execute(measurementID)
	if err != nil {
		return errorslist.ErrInternalMsg
	}
	return checkAccess(measurement.UserID, userID)
}
