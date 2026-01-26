package timers

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises"
	timerusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/timers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter *Presenter

	stopTimerUC  *timerusecases.StopUseCase
	startTimerUC *timerusecases.StartUseCase

	exerciseHandler *exercises.Handler
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	stopTimerUC *timerusecases.StopUseCase,
	startTimerUC *timerusecases.StartUseCase,
	exerciseHandler *exercises.Handler,
) *Handler {
	return &Handler{
		presenter: NewPresenter(bot),

		stopTimerUC:  stopTimerUC,
		startTimerUC: startTimerUC,

		exerciseHandler: exerciseHandler,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "timer_unpin_and_cancel_"):
		timerID := strings.TrimPrefix(data, "timer_unpin_and_cancel_")
		h.StopTimer(chatID, timerID)

	case strings.HasPrefix(data, "timer_start_"):
		parts := strings.Split(data, "_")
		if len(parts) >= 5 && parts[3] == "ex" {
			seconds, _ := strconv.Atoi(parts[2])
			exerciseID, _ := strconv.ParseInt(parts[4], 10, 64)
			h.StartTimer(chatID, exerciseID, seconds)
		}
	}
}

func (h *Handler) StopTimer(chatID int64, timerID string) {
	h.stopTimerUC.Execute(chatID, timerID)
}

func (h *Handler) StartTimer(chatID, exerciseID int64, seconds int) {
	res, err := h.startTimerUC.Execute(chatID, exerciseID, seconds)
	if err != nil {
		if errors.Is(err, timerusecases.TimerNotSupported) {
			h.presenter.ShowTimerIsNotSupported(chatID)
		}
		return
	}
	doWhenTimerExpired := func() { h.exerciseHandler.ShowCurrentExerciseSession(chatID, res.Exercise.WorkoutDayID) }
	h.presenter.ShowCreatedTimer(chatID, seconds, res, doWhenTimerExpired)
}
