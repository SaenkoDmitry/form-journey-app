package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) PushSubscribe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PushSubscribe")
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sub dto.PushSubscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.container.CreatePushSubscriptionUC.Execute(claims.UserID, sub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *serviceImpl) PushUnsubscribe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PushUnsubscribe")
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sub dto.PushUnsubscribe
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.container.DeletePushSubscriptionUC.Execute(claims.UserID, sub.Endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
