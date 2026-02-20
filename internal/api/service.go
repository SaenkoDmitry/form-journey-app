package api

import (
	"net/http"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/push"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/timermanager"
	"gorm.io/gorm"
)

type Service interface {
	TelegramRedirectHandler(w http.ResponseWriter, r *http.Request)
	TelegramLoginHandler(w http.ResponseWriter, r *http.Request)
	GetIcon(w http.ResponseWriter, r *http.Request)
	ChangeIcon(w http.ResponseWriter, r *http.Request)
	GetAllWorkouts(w http.ResponseWriter, r *http.Request)
	StartWorkout(w http.ResponseWriter, r *http.Request)
	FinishWorkout(w http.ResponseWriter, r *http.Request)
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
	ShowCurrentExerciseSession(w http.ResponseWriter, r *http.Request)
	MoveToExerciseSession(w http.ResponseWriter, r *http.Request)
	AddSet(w http.ResponseWriter, r *http.Request)
	DeleteSet(w http.ResponseWriter, r *http.Request)
	CompleteSet(w http.ResponseWriter, r *http.Request)
	ChangeSet(w http.ResponseWriter, r *http.Request)
	DeleteExercise(w http.ResponseWriter, r *http.Request)
	AddExercise(w http.ResponseWriter, r *http.Request)
	PushSubscribe(w http.ResponseWriter, r *http.Request)
	PushUnsubscribe(w http.ResponseWriter, r *http.Request)
	StartTimer(w http.ResponseWriter, r *http.Request)
	CancelTimer(w http.ResponseWriter, r *http.Request)
}

type serviceImpl struct {
	container    *usecase.Container
	timerManager *timermanager.TimerManager
}

func New(container *usecase.Container, db *gorm.DB) Service {
	pushService := push.NewService(db)
	return &serviceImpl{
		container:    container,
		timerManager: timermanager.NewTimerManager(db, pushService),
	}
}
