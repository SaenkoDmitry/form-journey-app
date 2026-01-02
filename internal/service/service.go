package service

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service interface {
	HandleMessage(message *tgbotapi.Message)
	HandleCallback(callback *tgbotapi.CallbackQuery)
}

type serviceImpl struct {
	bot           *tgbotapi.BotAPI
	usersRepo     users.Repo
	workoutsRepo  workouts.Repo
	exercisesRepo exercises.Repo
	setsRepo      sets.Repo
	sessionsRepo  sessions.Repo
}

func NewService(
	bot *tgbotapi.BotAPI,
	usersRepo users.Repo,
	workoutsRepo workouts.Repo,
	exercisesRepo exercises.Repo,
	setsRepo sets.Repo,
	sessionsRepo sessions.Repo,
) Service {
	return &serviceImpl{
		bot:           bot,
		usersRepo:     usersRepo,
		workoutsRepo:  workoutsRepo,
		exercisesRepo: exercisesRepo,
		setsRepo:      setsRepo,
		sessionsRepo:  sessionsRepo,
	}
}
