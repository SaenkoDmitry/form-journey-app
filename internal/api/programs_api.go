package api

import (
	"encoding/json"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/validator"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) GetUserPrograms(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	result, err := s.container.FindAllProgramsByUserUC.Execute(claims.UserID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Programs)
}

func (s *serviceImpl) GetActiveProgramForUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.container.GetUserByIDUC.Execute(claims.UserID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if user.ActiveProgramID == nil {
		http.Error(w, "У вас нет активных программ, создайте хотя бы одну", http.StatusForbidden)
		return
	}

	program, err := s.container.GetProgramUC.Execute(*user.ActiveProgramID, claims.UserID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

func (s *serviceImpl) CreateProgram(w http.ResponseWriter, r *http.Request) {
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

	err := s.container.CreateProgramUC.Execute(claims.UserID, input.Name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) ChooseProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err = s.container.ActivateProgramUC.Execute(claims.UserID, programID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	err = s.container.DeleteProgramUC.Execute(claims.UserID, programID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) RenameProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
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

	err = s.container.RenameProgramUC.Execute(programID, input.Name)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}

func (s *serviceImpl) GetProgram(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = validator.ValidateAccessToProgram(s.container, claims.UserID, programID); err != nil {
		helpers.WriteError(w, err)
		return
	}

	program, err := s.container.GetProgramUC.Execute(programID, claims.UserID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}
