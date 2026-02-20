package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

func (s *serviceImpl) GetMeasurements(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	offset, limit := helpers.GetOffsetLimit(r, 10, 50)

	result, err := s.container.FindAllMeasurementsUC.Execute(claims.UserID, limit, offset)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *serviceImpl) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	measurementID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err := validator.ValidateAccessToMeasurement(s.container, claims.UserID, measurementID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err := s.container.DeleteMeasurementByIDUC.Execute(measurementID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		Shoulders int `json:"shoulders"`
		Chest     int `json:"chest"`
		HandLeft  int `json:"hand_left"`
		HandRight int `json:"hand_right"`
		Waist     int `json:"waist"`
		Buttocks  int `json:"buttocks"`
		HipLeft   int `json:"hip_left"`
		HipRight  int `json:"hip_right"`
		CalfLeft  int `json:"calf_left"`
		CalfRight int `json:"calf_right"`
		Weight    int `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Создаём модель
	m := &models.Measurement{
		UserID:    claims.UserID,
		CreatedAt: time.Now(),
		Shoulders: input.Shoulders * 10,
		Chest:     input.Chest * 10,
		HandLeft:  input.HandLeft * 10,
		HandRight: input.HandRight * 10,
		Waist:     input.Waist * 10,
		Buttocks:  input.Buttocks * 10,
		HipLeft:   input.HipLeft * 10,
		HipRight:  input.HipRight * 10,
		CalfLeft:  input.CalfLeft * 10,
		CalfRight: input.CalfRight * 10,
		Weight:    input.Weight * 1000,
	}

	result, err := s.container.CreateMeasurementUC.Execute(m)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
