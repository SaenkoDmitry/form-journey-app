package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) ShowCurrentExerciseSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutID, err := strconv.ParseInt(r.PathValue("workout_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	session, err := s.container.ShowCurrentExerciseSessionUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (s *serviceImpl) MoveToExerciseSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutID, err := strconv.ParseInt(r.PathValue("workout_id"), 10, 64)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		Next bool `json:"next"` // Если false, то двигаемся к предыдущему упражнению, если true – к следующему
	}

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err = validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err = s.container.MoveSessionToExerciseUC.Execute(workoutID, input.Next)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) MoveToCertainExerciseSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutID, err := helpers.ParseInt64Param("workout_id", w, r)
	if err != nil {
		return
	}

	index, err := helpers.ParseInt64Param("index", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err = s.container.MoveToCertainUC.Execute(workoutID, int(index))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
