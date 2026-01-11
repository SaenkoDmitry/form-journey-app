package models

import "time"

type User struct {
	ID              int64 `gorm:"primaryKey;autoIncrement"`
	Username        string
	ChatID          int64
	FirstName       string
	LastName        string
	LanguageCode    string
	ActiveProgramID *int64
	CreatedAt       time.Time
	Programs        []WorkoutProgram `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (*User) TableName() string {
	return "users"
}
