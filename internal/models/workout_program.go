package models

import "time"

type WorkoutProgram struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64
	Name      string
	CreatedAt time.Time
	DayTypes  []WorkoutDayType `gorm:"foreignKey:WorkoutProgramID;constraint:OnDelete:CASCADE"`
}

func (*WorkoutProgram) TableName() string {
	return "workout_programs"
}
