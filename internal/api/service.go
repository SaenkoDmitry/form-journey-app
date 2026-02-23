package api

import (
	"net/http"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/push"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/timermanager"
	"gorm.io/gorm"
)

type Service interface {
	MeHandler(w http.ResponseWriter, r *http.Request)

	// ----- telegram auth -----

	TelegramRedirectHandler(w http.ResponseWriter, r *http.Request)
	TelegramLoginHandler(w http.ResponseWriter, r *http.Request)

	// ----- yandex auth -----

	YandexRedirectHandler(w http.ResponseWriter, r *http.Request)
	YandexLoginHandler(w http.ResponseWriter, r *http.Request)

	// ----- user profile icon -----

	GetIcon(w http.ResponseWriter, r *http.Request)
	ChangeIcon(w http.ResponseWriter, r *http.Request)

	// ----- workouts -----

	GetAllWorkouts(w http.ResponseWriter, r *http.Request)
	StartWorkout(w http.ResponseWriter, r *http.Request)
	FinishWorkout(w http.ResponseWriter, r *http.Request)
	ReadWorkout(w http.ResponseWriter, r *http.Request)
	DeleteWorkout(w http.ResponseWriter, r *http.Request)

	// ----- measurements -----

	GetMeasurements(w http.ResponseWriter, r *http.Request)
	CreateMeasurement(w http.ResponseWriter, r *http.Request)
	DeleteMeasurement(w http.ResponseWriter, r *http.Request)

	// ----- exercise group types -----

	GetExerciseGroups(w http.ResponseWriter, _ *http.Request)
	GetExerciseTypesByGroup(w http.ResponseWriter, r *http.Request)

	// ----- programs -----

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

	// ----- presets -----

	ParsePreset(w http.ResponseWriter, r *http.Request)
	SavePreset(w http.ResponseWriter, r *http.Request)

	// ----- sessions -----

	ShowCurrentExerciseSession(w http.ResponseWriter, r *http.Request)
	MoveToExerciseSession(w http.ResponseWriter, r *http.Request)

	// ----- sets -----

	AddSet(w http.ResponseWriter, r *http.Request)
	DeleteSet(w http.ResponseWriter, r *http.Request)
	CompleteSet(w http.ResponseWriter, r *http.Request)
	ChangeSet(w http.ResponseWriter, r *http.Request)

	// ----- exercises -----

	DeleteExercise(w http.ResponseWriter, r *http.Request)
	AddExercise(w http.ResponseWriter, r *http.Request)

	// ----- notifications -----

	PushSubscribe(w http.ResponseWriter, r *http.Request)
	PushUnsubscribe(w http.ResponseWriter, r *http.Request)

	// ----- timers -----

	StartTimer(w http.ResponseWriter, r *http.Request)
	CancelTimer(w http.ResponseWriter, r *http.Request)

	// ----- video -----

	StreamVideo(w http.ResponseWriter, r *http.Request)
	LinkVideo(w http.ResponseWriter, r *http.Request)
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

func (s *serviceImpl) isAllowedOrigin(origin string) bool {
	if strings.HasSuffix(origin, ".lhr.life") {
		return true
	}
	allowed := []string{
		"http://localhost:3000",
		"https://form-journey.ru",
		"https://96737811-dd90-496f-ac88-15158530e662-e1.tunnel4.com",
	}

	for _, o := range allowed {
		if o == origin {
			return true
		}
	}

	return false
}
