package dto

import "github.com/SaenkoDmitry/training-tg-bot/internal/models"

type DeleteProgramResult struct {
	User *models.User
}

type ActivateProgramResult struct {
	User *models.User
}

type GetAllPrograms struct {
	User     *models.User
	Programs []models.WorkoutProgram
}

type CreateProgramResult struct {
	User *models.User
}

type ListGroups struct {
	Groups []models.ExerciseGroupType
}

type GetProgram struct {
	Program          models.WorkoutProgram
	ExerciseTypesMap map[int64]models.ExerciseType
}

type ConfirmDeleteProgram struct {
	Program models.WorkoutProgram
}

type RenameProgram struct {
	Program models.WorkoutProgram
}
