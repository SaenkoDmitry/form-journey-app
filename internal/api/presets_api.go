package api

import (
	"encoding/json"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

func (s *serviceImpl) ParsePreset(w http.ResponseWriter, r *http.Request) {

	// Разбираем JSON из тела запроса
	var input struct {
		Preset string `json:"preset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	exercisesList, err := s.container.ExerciseTypeListUC.Execute()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	exercisesMap := make(map[int64]models.ExerciseType)
	for _, ex := range exercisesList.ExerciseTypes {
		exercisesMap[ex.ID] = ex
	}

	result := &dto.PresetListDTO{
		Exercises: make([]*dto.ExerciseDTO, 0),
	}

	exercisesPreset := utils.SplitPreset(input.Preset)
	if len(exercisesPreset) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}

	for _, ex := range exercisesPreset {
		sets := make([]*dto.SetDTO, 0, len(ex.Sets))
		for _, set := range ex.Sets {
			sets = append(sets, &dto.SetDTO{
				Reps:    set.Reps,
				Weight:  set.Weight,
				Minutes: set.Minutes,
			})
		}
		result.Exercises = append(result.Exercises, &dto.ExerciseDTO{
			ID:    ex.ID,
			Units: exercisesMap[ex.ID].Units,
			Name:  exercisesMap[ex.ID].Name,
			Sets:  sets,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *serviceImpl) SavePreset(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		DayTypeID int64  `json:"day_type_id"`
		NewPreset string `json:"new_preset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	day, err := s.container.GetDayTypeUC.Execute(input.DayTypeID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, day.WorkoutProgramID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err = s.container.UpdatePresetUC.Execute(day.ID, input.NewPreset)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
