package telegram

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/changes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exports"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/stats"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/timers"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/router"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type App struct {
	bot    *tgbotapi.BotAPI
	router *router.Router
}

func New(token string, useCases *usecase.Container) (*App, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	workoutsHandler := workouts.NewHandler(
		bot,
		useCases.DeleteWorkoutUC,
		useCases.ConfirmDeleteWorkoutUC,
		useCases.CreateWorkoutUC,
		useCases.StartWorkoutUC,
		useCases.FindMyWorkoutsUC,
		useCases.ShowWorkoutProgressUC,
		useCases.FinishWorkoutUC,
		useCases.ConfirmFinishWorkoutUC,
		useCases.ShowCurrentExerciseSessionUC,
		useCases.FindWorkoutsByUserUC,
		useCases.StatsWorkoutUC,
		useCases.GetByUserProgramUC,
	)

	exercisesHandler := exercises.NewHandler(
		bot,
		useCases.ShowCurrentExerciseSessionUC,
		useCases.GetGroupUC,
		useCases.FindTypesByGroupUC,
		useCases.ConfirmDeleteExerciseUC,
		useCases.DeleteExerciseUC,
		useCases.MoveSessionToExerciseUC,
		useCases.GetExerciseUC,
		useCases.GetAllGroupsUC,
		useCases.CreateExerciseUC,
		workoutsHandler,
	)

	timersHandler := timers.NewHandler(bot, useCases.StopTimerUC, useCases.StartTimerUC, exercisesHandler)

	statsHandler := stats.NewHandler(bot, useCases.PeriodStatsUC)

	setsHandler := sets.NewHandler(bot,
		useCases.CompleteSetUC, useCases.AddOneMoreSetUC, useCases.RemoveLastSetUC,
		useCases.ShowCurrentExerciseSessionUC, exercisesHandler, timersHandler,
	)

	programsHandler := programs.NewHandler(
		bot,
		useCases.DeleteProgramUC,
		useCases.CreateProgramUC,
		useCases.ActivateProgramUC,
		useCases.GetProgramUC,
		useCases.FindAllProgramsByUserUC,
	)

	usersHandler := users.NewHandler(bot, useCases.CreateUserUC)

	dayTypesHandler := daytypes.NewHandler(bot, useCases.GetAllGroupsUC)

	exportsHandler := exports.NewHandler(bot, useCases.ExportToExcelUC)

	changesHandler := changes.NewHandler(bot,
		useCases.ShowCurrentExerciseSessionUC, useCases.UpdateNextSetUC, useCases.FindAllProgramsByUserUC, useCases.RenameProgramUC,
		useCases.GetAllGroupsUC, useCases.DayTypesCreateUC, useCases.UpdateDateTypeUC, useCases.GetDayTypeUC,
		useCases.ExerciseTypeListUC, useCases.GetProgramUC)

	r := router.New(
		bot,
		useCases.CreateUserUC,
		useCases.GetUserUC,
		usersHandler,
		workoutsHandler,
		timersHandler,
		statsHandler,
		setsHandler,
		programsHandler,
		exportsHandler,
		exercisesHandler,
		changesHandler,
		dayTypesHandler,
	)

	return &App{
		bot:    bot,
		router: r,
	}, nil
}

func (a *App) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)

	for update := range updates {
		a.router.HandleUpdate(update)
	}

	return nil
}
