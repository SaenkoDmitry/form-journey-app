package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) CreateProgramDay(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	dayTypeID, err := s.container.DayTypesCreateUC.Execute(programID, input.Name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	day, err := s.container.GetDayTypeUC.Execute(dayTypeID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(day)
}

func (s *serviceImpl) UpdateProgramDay(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	dayTypeID, err := helpers.ParseInt64Param("day_type_id", w, r)
	if err != nil {
		return
	}

	type inputSet struct {
		Reps    int     `json:"reps"`
		Weight  float32 `json:"weight"`
		Meters  int     `json:"meters"`
		Minutes int     `json:"minutes"`
	}

	// Разбираем JSON из тела запроса
	var input struct {
		ExerciseTypeID int64       `json:"exercise_type_id"`
		Sets           []*inputSet `json:"sets"`
	}

	formatSets := func(sets []*inputSet) string {
		arrayOfSets := make([]string, 0, len(sets))
		for _, set := range sets {
			switch {
			case set.Weight > 0:
				arrayOfSets = append(arrayOfSets, fmt.Sprintf("%d*%.0f", set.Reps, set.Weight))
			case set.Minutes > 0:
				arrayOfSets = append(arrayOfSets, fmt.Sprintf("%d", set.Minutes))
			case set.Meters > 0:
				arrayOfSets = append(arrayOfSets, fmt.Sprintf("%d", set.Meters))
			}
		}
		return strings.Join(arrayOfSets, ";")
	}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = s.container.AddExPresetUC.Execute(dayTypeID, input.ExerciseTypeID, formatSets(input.Sets))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) GetProgramDay(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	dayTypeID, err := helpers.ParseInt64Param("day_type_id", w, r)
	if err != nil {
		return
	}

	day, err := s.container.GetDayTypeUC.Execute(dayTypeID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(day)
}

func (s *serviceImpl) DeleteProgramDay(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	dayTypeID, err := helpers.ParseInt64Param("day_type_id", w, r)
	if err != nil {
		return
	}

	err = s.container.DeleteDayTypeUC.Execute(dayTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
