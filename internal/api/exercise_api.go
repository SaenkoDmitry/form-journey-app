package api

import (
	"encoding/json"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exerciseID, err := helpers.ParseInt64Param("id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToExercise(s.container, claims.UserID, exerciseID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	_, err = s.container.DeleteExerciseUC.Execute(exerciseID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) AddExercise(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		WorkoutID      int64 `json:"workout_id"`
		ExerciseTypeID int64 `json:"exercise_type_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateAccessToWorkout(s.container, claims.UserID, input.WorkoutID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	if _, err := s.container.CreateExerciseUC.Execute(input.WorkoutID, input.ExerciseTypeID); err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
