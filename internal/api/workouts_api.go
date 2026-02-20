package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	offset, limit := helpers.GetOffsetLimit(r, 10, 50)

	res, err := s.container.FindMyWorkoutsUC.Execute(claims.UserID, offset, limit)
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
		w.WriteHeader(http.StatusUnauthorized)
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

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, day.WorkoutProgramID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	createdWorkout, err := s.container.CreateWorkoutUC.Execute(claims.UserID, input.DayTypeID) // создаем тренировку
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutIDStr := r.PathValue("workout_id")
	workoutID, _ := strconv.ParseInt(workoutIDStr, 10, 64)

	if err := validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	progress, err := s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutID, err := helpers.ParseInt64Param("workout_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	workoutID, err := helpers.ParseInt64Param("workout_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToWorkout(s.container, claims.UserID, workoutID); err != nil {
		helpers.WriteError(w, err)
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
