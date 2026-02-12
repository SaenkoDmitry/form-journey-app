package api

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"net/http"
)

type Service interface {
	GetAllWorkouts(w http.ResponseWriter, r *http.Request)
	ReadWorkout(w http.ResponseWriter, r *http.Request)
	DeleteWorkout(w http.ResponseWriter, r *http.Request)
	GetMeasurements(w http.ResponseWriter, r *http.Request)
	CreateMeasurement(w http.ResponseWriter, r *http.Request)
	DeleteMeasurement(w http.ResponseWriter, r *http.Request)
	GetExerciseGroups(w http.ResponseWriter, _ *http.Request)
	GetExerciseTypesByGroup(w http.ResponseWriter, r *http.Request)
	GetUserPrograms(w http.ResponseWriter, r *http.Request)
	GetActiveProgramForUser(w http.ResponseWriter, r *http.Request)
	CreateProgram(w http.ResponseWriter, r *http.Request)
	ChooseProgram(w http.ResponseWriter, r *http.Request)
	DeleteProgram(w http.ResponseWriter, r *http.Request)
	RenameProgram(w http.ResponseWriter, r *http.Request)
	GetProgram(w http.ResponseWriter, r *http.Request)
	CreateProgramDay(w http.ResponseWriter, r *http.Request)
	DeleteProgramDay(w http.ResponseWriter, r *http.Request)
	UpdateProgramDay(w http.ResponseWriter, r *http.Request)
	GetProgramDay(w http.ResponseWriter, r *http.Request)
	ParsePreset(w http.ResponseWriter, r *http.Request)
	SavePreset(w http.ResponseWriter, r *http.Request)
}

type serviceImpl struct {
	container *usecase.Container
}

func New(container *usecase.Container) Service {
	return &serviceImpl{
		container: container,
	}
}
