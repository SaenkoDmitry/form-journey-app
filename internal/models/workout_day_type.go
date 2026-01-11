package models

import "time"

type WorkoutDayType struct {
	ID               int64 `gorm:"primaryKey;autoIncrement"`
	WorkoutProgramID int64
	Name             string
	Preset           string
	CreatedAt        time.Time
}

func (*WorkoutDayType) TableName() string {
	return "workout_day_types"
}
