package models

import "time"

type WorkoutShare struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	WorkoutDayID int64     `gorm:"not null;index"`
	Token        string    `gorm:"type:varchar(64);not null;uniqueIndex"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	ExpiresAt    *time.Time
}

func (*WorkoutShare) TableName() string {
	return "workout_shares"
}
