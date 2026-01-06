package utils

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
)

func GetWorkoutNameByID(ID string) string {
	switch ID {
	case constants.LegsAndShouldersWorkoutID:
		return constants.LegsAndShouldersWorkoutName
	case constants.BackAndBicepsWorkoutID:
		return constants.BackAndBicepsWorkoutName
	case constants.ChestAndTricepsID:
		return constants.ChestAndTricepsName
	case constants.CardioID:
		return constants.CardioName
	}
	return ""
}
