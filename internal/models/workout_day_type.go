package models

import "time"

type WorkoutDayType struct {
	ID               int64 `gorm:"primaryKey"`
	WorkoutProgramID int64
	Name             string
	Preset           string
	CreatedAt        time.Time
}

func (*WorkoutDayType) TableName() string {
	return "training.workout_day_types"
}
