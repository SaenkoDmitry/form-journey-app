package dto

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"time"
)

type WorkoutItem struct {
	ID        int64
	Name      string
	StartedAt time.Time
	EndedAt   *time.Time
	Duration  string
	Completed bool
}

type Pagination struct {
	Offset int
	Limit  int
	Total  int
}

type ShowMyWorkoutsResult struct {
	Items      []WorkoutItem
	Pagination Pagination
}

type ConfirmDeleteWorkout struct {
	WorkoutID   int64
	DayTypeName string
}

type DeleteWorkout struct {
}

type ConfirmFinishWorkout struct {
	DayType models.WorkoutDayType
}

type FinishWorkout struct {
	WorkoutID int64
}

type CreateWorkout struct {
	WorkoutID int64
}

type StartWorkout struct {
}

type WorkoutProgress struct {
	Workout *models.WorkoutDay

	TotalExercises     int
	CompletedExercises int

	TotalSets     int
	CompletedSets int

	ProgressPercent int
	RemainingMin    *int

	SessionStarted bool
}

type WorkoutStatistic struct {
	DayType            models.WorkoutDayType
	WorkoutDay         models.WorkoutDay
	TotalWeight        float64
	CompletedExercises int
	TotalTime          int
	ExerciseTypesMap   map[int64]models.ExerciseType
	ExerciseWeightMap  map[int64]float64
	ExerciseTimeMap    map[int64]int
}

type ShowWorkoutByUserID struct {
	Workouts []models.WorkoutDay
	User     *models.User
}
