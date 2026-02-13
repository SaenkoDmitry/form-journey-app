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

func (s *serviceImpl) StartWorkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		DayTypeID int64 `json:"day_type_id"`
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

	program, err := s.container.GetProgramUC.Execute(day.WorkoutProgramID, claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if program.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	createdWorkout, err := s.container.CreateWorkoutUC.Execute(claims.ChatID, input.DayTypeID) // создаем тренировку
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	_, err = s.container.StartWorkoutUC.Execute(createdWorkout.WorkoutID) // создаем сессию
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&StartWorkoutDTO{WorkoutID: createdWorkout.WorkoutID})
}

type StartWorkoutDTO struct {
	WorkoutID int64 `json:"workout_id"`
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

func (s *serviceImpl) FinishWorkout(w http.ResponseWriter, r *http.Request) {
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

	_, err = s.container.FinishWorkoutUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
