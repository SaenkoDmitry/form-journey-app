package api

import (
	"encoding/json"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"net/http"
	"strconv"
)

func (s *serviceImpl) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	offset, limit := helpers.GetOffsetLimit(r, 10, 50)

	res, err := s.container.FindMyWorkoutsUC.Execute(claims.ChatID, offset, limit)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *serviceImpl) ReadWorkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workoutIDStr := r.PathValue("workout_id")
	workoutID, _ := strconv.ParseInt(workoutIDStr, 10, 64)

	progress, err := s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if progress.Workout.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	stats, err := s.container.StatsWorkoutUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&ReadWorkoutDTO{Progress: progress, Stats: stats})
}

type ReadWorkoutDTO struct {
	Progress *dto.WorkoutProgress `json:"progress"`
	Stats    *dto.WorkoutStatistic
}

func (s *serviceImpl) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workoutIDStr := r.PathValue("workout_id")
	workoutID, err := strconv.ParseInt(workoutIDStr, 10, 64)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	progress, err := s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if progress.Workout.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	err = s.container.DeleteWorkoutUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
