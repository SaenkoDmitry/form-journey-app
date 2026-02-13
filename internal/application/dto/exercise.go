package dto

import "github.com/SaenkoDmitry/training-tg-bot/internal/models"

type CurrentExerciseSession struct {
	Exercise      *FormattedExercise `json:"exercise"`
	ExerciseObj   *ExerciseTypeDTO   `json:"exercise_type"`
	DayType       *WorkoutDayTypeDTO `json:"day_type"`
	WorkoutDay    *FormattedWorkout  `json:"workout"`
	ExerciseIndex int                `json:"exercise_index"`
}

type ExerciseTypeList struct {
	ExerciseTypes []models.ExerciseType
}

type FindTypesByGroup struct {
	ExerciseTypes []*ExerciseTypeDTO `json:"exercise_types"`
}

type ExerciseTypeDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Url           string `json:"url"`
	GroupName     string `json:"group_name"`
	RestInSeconds int    `json:"rest_in_seconds"`
	Accent        string `json:"accent"`
	Units         string `json:"units"`
	Description   string `json:"description"`
}

func MapExerciseTypeDTOList(types []models.ExerciseType, groupsMap map[string]string) []*ExerciseTypeDTO {
	result := make([]*ExerciseTypeDTO, 0, len(types))
	for _, t := range types {
		result = append(result, MapExerciseTypeDTO(t, groupsMap))
	}
	return result
}

func MapExerciseTypeDTO(t models.ExerciseType, groupsMap map[string]string) *ExerciseTypeDTO {
	return &ExerciseTypeDTO{
		ID:            t.ID,
		Name:          t.Name,
		Url:           t.Url,
		GroupName:     groupsMap[t.ExerciseGroupTypeCode],
		RestInSeconds: t.RestInSeconds,
		Accent:        t.Accent,
		Units:         t.Units,
		Description:   t.Description,
	}
}

type ConfirmDeleteExercise struct {
	Exercise    models.Exercise
	ExerciseObj models.ExerciseType
}

type GetExercise struct {
	ExerciseType models.ExerciseType
}

type CreateExercise struct {
	ExerciseObj models.ExerciseType
}
