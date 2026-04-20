package models

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID int64 `gorm:"primaryKey;autoIncrement"`

	Username string // telegram
	ChatID   int64  // telegram

	FirstName       string
	LastName        string
	Email           string
	LanguageCode    string
	Icon            string
	ActiveProgramID *int64
	CreatedAt       time.Time
	Programs        []WorkoutProgram `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	YandexID    string // yandex
	YandexLogin string // yandex
}

func (u *User) GetFirstName() string {
	if u == nil {
		return ""
	}
	return u.FirstName
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) IsAdmin() bool {
	return u.ID == 1 && u.Username == "dsaenko"
}

func (u *User) FullName() string {
	arr := make([]string, 0)
	if u.FirstName != "" {
		arr = append(arr, u.FirstName)
	}
	if u.LastName != "" {
		arr = append(arr, u.LastName)
	}
	return fmt.Sprintf("%s (%s)", strings.Join(arr, " "), u.Username)
}

func (u *User) ShortName() string {
	arr := make([]string, 0)
	if u.FirstName != "" {
		arr = append(arr, u.FirstName)
	}
	if u.LastName != "" {
		arr = append(arr, u.LastName)
	}
	return fmt.Sprintf("%s", strings.Join(arr, " "))
}
