package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram"
	"github.com/SaenkoDmitry/training-tg-bot/internal/api"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	"github.com/SaenkoDmitry/training-tg-bot/internal/middlewares"
	"github.com/SaenkoDmitry/training-tg-bot/internal/web"
)

func main() {
	log.Println("DEPLOY MARK:", time.Now().UnixNano())

	token := os.Getenv("TELEGRAM_TOKEN")
	dsn := os.Getenv("DATABASE_URL")
	fmt.Printf("TELEGRAM_TOKEN: %s, DATABASE_URL: %s\n", token, dsn)

	// init database
	db := initDB(dsn)

	// health check handler
	go healthCheckAPIHandler(db)

	// use cases
	container := usecase.NewContainer(db)

	// init telegram app
	app, err := telegram.New(token, container)
	if err != nil {
		return
	}

	go initServer(container)

	// run telegram app
	if err = app.Run(); err != nil {
		log.Panic(err)
	}
}

func initDB(dsn string) *gorm.DB {
	migrate(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Failed to connect database: %s", err.Error())
	}
	return db
}

func healthCheckAPIHandler(db *gorm.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	dbObj, err := db.DB()
	if err != nil {
		log.Fatalf("cannot get *sql.DB")
	}
	fmt.Println("port:", port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err = dbObj.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "DB ping failed")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Printf("Health check server failed: %s", err)
	}
}

func migrate(dsn string) {
	migrateDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Panicf("Failed to migrate: open dsn: %s", err.Error())
	}
	defer migrateDB.Close()

	if err = goose.SetDialect("postgres"); err != nil {
		log.Panicf("Failed to migrate: set database dialect: %s", err.Error())
	}

	if err := goose.Up(migrateDB, "db/migrations"); err != nil {
		log.Panicf("Failed to migrate: goose up: %s", err.Error())
	}
}

func initServer(container *usecase.Container) {
	r := chi.NewRouter()

	// Средства middleware chi
	r.Use(middleware.Logger)    // лог запросов
	r.Use(middleware.Recoverer) // recovery от паник
	r.Use(middleware.RequestID) // уникальный ID запроса

	s := api.New(container)

	r.Route("/api/login", func(r chi.Router) {
		r.Post("/", api.LoginHandler)
	})
	r.Route("/api/logout", func(r chi.Router) {
		r.Post("/", api.LogoutHandler)
	})

	r.Route("/api/me", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", api.MeHandler)
	})

	r.Route("/api/workouts", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetAllWorkouts)          // GET /api/workouts
		r.Get("/{workout_id}", s.ReadWorkout) // GET /api/workouts/123
	})

	r.Route("/api/measurements", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetMeasurements)
		r.Post("/", s.CreateMeasurement)
	})

	r.Route("/api/exercise-groups", func(r chi.Router) {
		r.Use(middlewares.Auth)

		r.Get("/", s.GetExerciseGroups)
		r.Get("/{group}", s.GetExerciseTypesByGroup)
	})

	// UI (React build)
	r.NotFound(web.SPAHandler().ServeHTTP)

	log.Println("Server started on :8080")
	http.ListenAndServe("0.0.0.0:8080", r)
}
