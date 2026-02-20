package daytypes

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/daytypes"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/groups"
	programusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	getAllGroupsUC *groups.GetAllUseCase
	getDayUC       *daytypes.GetUseCase
	deleteUC       *daytypes.DeleteUseCase
	getProgramUC   *programusecases.GetUseCase

	programsHandler *programs.Handler
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	getUC *daytypes.GetUseCase,
	getAllGroupsUC *groups.GetAllUseCase,
	deleteUC *daytypes.DeleteUseCase,
	getProgramUC *programusecases.GetUseCase,
	programsHandler *programs.Handler,
) *Handler {
	return &Handler{
		presenter:       NewPresenter(bot),
		commonPresenter: common.NewPresenter(bot),
		getDayUC:        getUC,
		getAllGroupsUC:  getAllGroupsUC,
		deleteUC:        deleteUC,
		getProgramUC:    getProgramUC,
		programsHandler: programsHandler,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "day_type_view_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "day_type_view_"), 10, 64)
		h.ViewDayType(chatID, dayTypeID)

	case strings.HasPrefix(data, "day_type_edit_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "day_type_edit_"), 10, 64)
		h.selectDayTypeExercise(chatID, dayTypeID)

	case strings.HasPrefix(data, "day_type_confirm_delete_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "day_type_confirm_delete_"), 10, 64)
		h.confirmDeleteDayType(chatID, dayTypeID)

	case strings.HasPrefix(data, "day_type_delete_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "day_type_delete_"), 10, 64)
		h.deleteDayType(chatID, dayTypeID)
	}
}

func (h *Handler) selectDayTypeExercise(chatID int64, dayTypeID int64) {
	res, err := h.getAllGroupsUC.Execute()
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getAllGroupsUC.Name())
		return
	}
	h.presenter.ShowSelectDayTypeDialog(chatID, dayTypeID, res)
}

func (h *Handler) deleteDayType(chatID int64, dayTypeID int64) {
	dayResult, err := h.getDayUC.Execute(dayTypeID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getDayUC.Name())
		return
	}
	programResult, err := h.getProgramUC.ExecuteByChatID(dayResult.WorkoutProgramID, chatID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getProgramUC.Name())
		return
	}
	err = h.deleteUC.Execute(dayTypeID)
	if err != nil {
		if errors.Is(err, daytypes.CannotDeleteAlreadyUsedDay) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.CannotDeleteDayTypeAlreadyUsedInWorkoutDays)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.deleteUC.Name())
		return
	}
	h.programsHandler.ViewProgram(chatID, programResult.ID)
}

func (h *Handler) confirmDeleteDayType(chatID int64, dayTypeID int64) {
	res, err := h.getDayUC.Execute(dayTypeID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getDayUC.Name())
		return
	}

	h.presenter.ShowConfirmDelete(chatID, res)
}

func (h *Handler) ViewDayType(chatID int64, dayTypeID int64) {
	day, err := h.getDayUC.Execute(dayTypeID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getDayUC.Name())
		return
	}
	program, err := h.getProgramUC.ExecuteByChatID(day.WorkoutProgramID, chatID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getDayUC.Name())
		return
	}
	h.presenter.ViewDayType(chatID, day, program)
}
