package api

import (
	"encoding/json"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) GetIcon(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.container.GetUserByIDUC.Execute(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Name string `json:"name"`
	}{
		Name: user.Icon,
	})
}

func (s *serviceImpl) ChangeIcon(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Разбираем JSON из тела запроса
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := s.container.ChangeIconUC.Execute(claims.UserID, input.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
