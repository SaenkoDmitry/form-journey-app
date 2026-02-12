package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/api/helpers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
)

func (s *serviceImpl) GetUserPrograms(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.FromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	result, err := s.container.FindAllProgramsByUserUC.Execute(claims.ChatID)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := s.container.GetUserUC.Execute(claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if user.ActiveProgramID == nil {
		http.Error(w, "У вас нет активных программ, создайте хотя бы одну", http.StatusForbidden)
		return
	}

	program, err := s.container.GetProgramUC.Execute(*user.ActiveProgramID, claims.ChatID)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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

	err := s.container.CreateProgramUC.Execute(claims.ChatID, input.Name)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = s.validateAccessToProgram(w, claims.ChatID, programID); err != nil {
		return
	}

	err = s.container.ActivateProgramUC.Execute(claims.ChatID, programID)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = s.validateAccessToProgram(w, claims.ChatID, programID); err != nil {
		return
	}

	err = s.container.DeleteProgramUC.Execute(claims.ChatID, programID)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = s.validateAccessToProgram(w, claims.ChatID, programID); err != nil {
		return
	}

	newProgramName := r.PathValue("program_new_name")

	err = s.container.RenameProgramUC.Execute(programID, newProgramName)
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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	programID, err := helpers.ParseInt64Param("program_id", w, r)
	if err != nil {
		return
	}

	if err = s.validateAccessToProgram(w, claims.ChatID, programID); err != nil {
		return
	}

	program, err := s.container.GetProgramUC.Execute(programID, claims.ChatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(program)
}

func (s *serviceImpl) validateAccessToProgram(w http.ResponseWriter, chatID int64, programID int64) error {
	user, err := s.container.GetUserUC.Execute(chatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return err
	}

	program, err := s.container.GetProgramUC.Execute(programID, chatID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return err
	}

	// check access
	if program.UserID != user.ID {
		http.Error(w, "access denied", http.StatusForbidden)
		return errors.New("no access")
	}

	return nil
}
