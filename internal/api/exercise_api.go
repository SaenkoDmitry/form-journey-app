package api

import (
	"encoding/json"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"net/http"
)

func (s *serviceImpl) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	exerciseID, err := helpers.ParseInt64Param("id", w, r)
	if err != nil {
		return
	}

	err = s.validateAccessToExercise(w, claims.ChatID, exerciseID)
	if err != nil {
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

	err := s.validateAccessToWorkout(w, claims.ChatID, input.WorkoutID)
	if err != nil {
		return
	}

	_, err = s.container.CreateExerciseUC.Execute(input.WorkoutID, input.ExerciseTypeID)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
