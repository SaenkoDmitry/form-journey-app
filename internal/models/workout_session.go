package models

import "time"

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey;autoIncrement"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}

func (*WorkoutSession) TableName() string {
	return "workout_sessions"
}
