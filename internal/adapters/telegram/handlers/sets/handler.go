package sets

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises/presenter"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/timers"
	exercisecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/session"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	setusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/sets"
)

type Handler struct {
	presenter *Presenter

	commonPresenter   *common.Presenter
	exercisePresenter *presenter.Presenter

	completeSetUC   *setusecases.CompleteUseCase
	addOneMoreSetUC *setusecases.AddOneMoreUseCase
	removeLastSetUC *setusecases.RemoveLastUseCase

	showCurrentSessionUC *exercisecases.ShowCurrentExerciseSessionUseCase

	exerciseHandler *exercises.Handler
	timersHandler   *timers.Handler
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	completeSetUC *setusecases.CompleteUseCase,
	addOneMoreSetUC *setusecases.AddOneMoreUseCase,
	removeLastSetUC *setusecases.RemoveLastUseCase,
	showCurrentSessionUC *exercisecases.ShowCurrentExerciseSessionUseCase,
	exerciseHandler *exercises.Handler,
	timersHandler *timers.Handler,
) *Handler {
	return &Handler{
		presenter:            NewPresenter(bot),
		commonPresenter:      common.NewPresenter(bot),
		exercisePresenter:    presenter.NewPresenter(bot),
		completeSetUC:        completeSetUC,
		addOneMoreSetUC:      addOneMoreSetUC,
		removeLastSetUC:      removeLastSetUC,
		showCurrentSessionUC: showCurrentSessionUC,
		exerciseHandler:      exerciseHandler,
		timersHandler:        timersHandler,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "set_complete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_complete_"), 10, 64)
		h.completeExerciseSet(chatID, exerciseID)

	case strings.HasPrefix(data, "set_add_one_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_add_one_"), 10, 64)
		h.addOneMoreSet(chatID, exerciseID)

	case strings.HasPrefix(data, "set_remove_last_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "set_remove_last_"), 10, 64)
		h.removeLastSet(chatID, exerciseID)
	}
}

func (h *Handler) completeExerciseSet(chatID int64, exerciseID int64) {
	res, err := h.completeSetUC.Execute(exerciseID)
	if err != nil {
		if errors.Is(err, setusecases.DoNothingErr) {
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.completeSetUC.Name())
		return
	}

	h.commonPresenter.SendSimpleHtmlMessage(chatID, fmt.Sprintf(messages.SetCompleted))

	if res.NeedShowCurrent {
		if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(res.WorkoutID); sessionErr == nil {
			h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
		}
	}

	if res.NeedMoveToNext {
		h.exerciseHandler.MoveToNextExercise(chatID, res.WorkoutID)
	}

	// start timer
	if res.NeedStartTimer {
		h.timersHandler.StartTimer(chatID, exerciseID, res.Seconds)
	}
}

func (h *Handler) addOneMoreSet(chatID int64, exerciseID int64) {
	res, err := h.addOneMoreSetUC.Execute(exerciseID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.addOneMoreSetUC.Name())
		return
	}
	h.presenter.ShowOneMore(chatID)

	if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(res.WorkoutID); sessionErr == nil {
		h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
	}
}

func (h *Handler) removeLastSet(chatID int64, exerciseID int64) {
	res, err := h.removeLastSetUC.Execute(exerciseID)
	if err != nil {
		if errors.Is(err, setusecases.AddOneMoreExerciseToDeleteErr) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.AddOneMoreExerciseToDelete)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.addOneMoreSetUC.Name())
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.SetDeleted)

	if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(res.WorkoutID); sessionErr == nil {
		h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
	}
}
