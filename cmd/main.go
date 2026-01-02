package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/SaenkoDmitry/training-tg-bot/internal/config"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	userStates = make(map[int64]string)
)

func main() {
	configFile, err := os.Open("../config.json")
	if err != nil {
		log.Fatal("Config file not found")
	}
	defer configFile.Close()

	var config config.Config
	json.NewDecoder(configFile).Decode(&config)

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open("../workout_bot.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.WorkoutDay{}, &models.Exercise{}, &models.Set{}, &models.WorkoutSession{})

	log.Printf("Authorized on account %s", bot.Self.UserName)

	usersRepo := users.NewRepo(db)
	workoutsRepo := workouts.NewRepo(db)
	exercisesRepo := exercises.NewRepo(db)
	setsRepo := sets.NewRepo(db)
	sessionsRepo := sessions.NewRepo(db)

	svc := service.NewService(bot, usersRepo, workoutsRepo, exercisesRepo, setsRepo, sessionsRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			svc.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			svc.HandleCallback(update.CallbackQuery)
		}
	}
}
