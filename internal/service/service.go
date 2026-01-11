package service

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sync"
)

type Service interface {
	HandleMessage(message *tgbotapi.Message)
	HandleCallback(callback *tgbotapi.CallbackQuery)
}

type serviceImpl struct {
	bot *tgbotapi.BotAPI

	usersRepo              users.Repo
	programsRepo           programs.Repo
	dayTypesRepo           daytypes.Repo
	workoutsRepo           workouts.Repo
	exercisesRepo          exercises.Repo
	setsRepo               sets.Repo
	sessionsRepo           sessions.Repo
	exerciseTypesRepo      exercisetypes.Repo
	exerciseGroupTypesRepo exercisegrouptypes.Repo
	userStatesMachine      *UserStatesMachine
}

type UserStatesMachine struct {
	userStates map[int64]string
	mu         *sync.Mutex
}

func (u *UserStatesMachine) SetValue(key int64, value string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.userStates[key] = value
}

func (u *UserStatesMachine) GetValue(key int64) (string, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	v, ok := u.userStates[key]
	return v, ok
}

func NewUserStatesMachine() *UserStatesMachine {
	return &UserStatesMachine{
		userStates: make(map[int64]string),
		mu:         &sync.Mutex{},
	}
}

func NewService(
	bot *tgbotapi.BotAPI,
	usersRepo users.Repo,
	programsRepo programs.Repo,
	dayTypesRepo daytypes.Repo,
	workoutsRepo workouts.Repo,
	exercisesRepo exercises.Repo,
	exerciseTypesRepo exercisetypes.Repo,
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
	setsRepo sets.Repo,
	sessionsRepo sessions.Repo,
) Service {
	return &serviceImpl{
		bot:                    bot,
		usersRepo:              usersRepo,
		programsRepo:           programsRepo,
		workoutsRepo:           workoutsRepo,
		dayTypesRepo:           dayTypesRepo,
		exercisesRepo:          exercisesRepo,
		exerciseTypesRepo:      exerciseTypesRepo,
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
		setsRepo:               setsRepo,
		sessionsRepo:           sessionsRepo,
		userStatesMachine:      NewUserStatesMachine(),
	}
}
