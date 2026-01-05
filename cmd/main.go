package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sessions"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("telegram_token:", os.Getenv("telegram_token"))
	bot, err := tgbotapi.NewBotAPI(os.Getenv("telegram_token"))
	if err != nil {
		log.Panic(err)
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %s", err.Error())
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

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Bot is running")
		})
		log.Fatal(http.ListenAndServe(":5000", nil))
	}()

	for update := range updates {
		if update.Message != nil {
			svc.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			svc.HandleCallback(update.CallbackQuery)
		}
	}
}
