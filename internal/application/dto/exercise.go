package dto

import "github.com/SaenkoDmitry/training-tg-bot/internal/models"

type CurrentExerciseSession struct {
	Exercise      models.Exercise
	ExerciseObj   models.ExerciseType
	DayType       models.WorkoutDayType
	WorkoutDay    models.WorkoutDay
	ExerciseIndex int
}

type ExerciseTypeList struct {
	ExerciseTypes []models.ExerciseType
}

type FindTypesByGroup struct {
	ExerciseTypes []models.ExerciseType
}

type ConfirmDeleteExercise struct {
	Exercise    models.Exercise
	ExerciseObj models.ExerciseType
}

type GetExercise struct {
	Exercise models.Exercise
}

type CreateExercise struct {
	ExerciseObj models.ExerciseType
}
