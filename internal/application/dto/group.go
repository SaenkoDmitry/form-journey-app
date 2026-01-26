package dto

import "github.com/SaenkoDmitry/training-tg-bot/internal/models"

type Group struct {
	Code string
	Name string
}

type ExerciseGroupTypeList struct {
	Groups []models.ExerciseGroupType
}
