package models

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"strings"
)

type ExerciseType struct {
	ID                    int64 `gorm:"primaryKey;autoIncrement"`
	Name                  string
	Url                   string
	ExerciseGroupTypeCode string
	RestInSeconds         int
	Accent                string
	Units                 string
	Description           string
}

func (*ExerciseType) TableName() string {
	return "exercise_types"
}

func (t *ExerciseType) ShowReps() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.RepsUnit)
}

func (t *ExerciseType) ShowWeight() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.WeightUnit)
}

func (t *ExerciseType) ShowMinutes() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.MinutesUnit)
}

func (t *ExerciseType) ShowMeters() bool {
	if t == nil {
		return false
	}
	return strings.Contains(t.Units, constants.MetersUnit)
}
