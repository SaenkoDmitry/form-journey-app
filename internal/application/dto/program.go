package dto

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"time"
)

type DeleteProgramResult struct {
	User *models.User
}

type ActivateProgramResult struct {
	User *models.User
}

type GetAllPrograms struct {
	User     *models.User
	Programs []*ProgramDTO `json:"programs"`
}

type ProgramDTO struct {
	ID        int64                `json:"id"`
	UserID    int64                `json:"user_id"`
	Name      string               `json:"name"`
	CreatedAt string               `json:"created_at"`
	DayTypes  []*WorkoutDayTypeDTO `json:"day_types"`
	IsActive  bool                 `json:"is_active"`
}

func MapDayTypeDTO(obj models.WorkoutDayType) *WorkoutDayTypeDTO {
	return &WorkoutDayTypeDTO{
		ID:               obj.ID,
		WorkoutProgramID: obj.WorkoutProgramID,
		Name:             obj.Name,
		Preset:           obj.Preset,
		CreatedAt:        "📅 " + obj.CreatedAt.Add(time.Hour*3).Format("02.01.2006 15:04"),
	}
}

type WorkoutDayTypeDTO struct {
	ID               int64  `json:"id"`
	WorkoutProgramID int64  `json:"program_id"`
	Name             string `json:"name"`
	Preset           string `json:"preset"`
	CreatedAt        string `json:"created_at"`
}

type CreateProgramResult struct {
	User *models.User
}

type ListGroups struct {
	Groups []models.ExerciseGroupType
}

type GetProgramDTO struct {
	Program ProgramDTO
}

type ConfirmDeleteProgram struct {
	Program models.WorkoutProgram
}

type RenameProgram struct {
	Program models.WorkoutProgram
}
