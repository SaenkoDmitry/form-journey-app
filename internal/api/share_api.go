package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func (s *serviceImpl) CreateShareWorkout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rl, ok := middlewares.ShareLimiterFromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
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

	// Проверяем, что тренировка завершена
	var workoutStat *dto.WorkoutProgress
	workoutStat, err = s.container.ShowWorkoutProgressUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}
	if workoutStat == nil || !workoutStat.Workout.Completed {
		http.Error(w, "workout not completed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if shareModel, findShareErr := s.container.GetShareByWorkoutUC.Execute(workoutID); findShareErr == nil {
		json.NewEncoder(w).Encode(buildShareDTO(shareModel))
		return
	}

	if !rl.Allow(claims.UserID) {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	shareModel, err := s.container.CreateShareUC.Execute(workoutID)
	if err != nil {
		http.Error(w, "failed to create share", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(buildShareDTO(shareModel))
}

func buildShareDTO(shareModel *models.WorkoutShare) dto.ShareResponse {
	return dto.ShareResponse{
		Token:     shareModel.Token,
		ShareURL:  getShareURL(constants.Domain, shareModel.Token),
		CreatedAt: shareModel.CreatedAt.Add(3 * time.Hour).Format("02.01.2006 15:04"),
	}
}

func getShareURL(domain, token string) string {
	return fmt.Sprintf("%s/public/workouts/%s", domain, token)
}

func (s *serviceImpl) GetPublicWorkout(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		http.Error(w, "token required", http.StatusBadRequest)
		return
	}

	shareDTO, err := s.container.GetShareUC.Execute(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	workoutID := shareDTO.WorkoutDayID

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
	json.NewEncoder(w).Encode(&ReadWorkoutDTO{
		Progress:      progress,
		Stats:         stats,
		UserFirstName: stats.WorkoutDay.GetUser().GetFirstName(),
	})
}
