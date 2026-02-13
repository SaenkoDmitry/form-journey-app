package api

import (
	"encoding/json"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"net/http"
	"strconv"
)

func (s *serviceImpl) ShowCurrentExerciseSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := strconv.ParseInt(r.PathValue("workout_id"), 10, 64)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	workout, err := s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if workout.Workout.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		Next bool `json:"next"` // Если false, то двигаемся к предыдущему упражнению, если true – к следующему
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	workout, err := s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if workout.Workout.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
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
