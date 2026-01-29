package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/groups"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	getAllGroupsUC *groups.GetAllUseCase
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	getAllGroupsUC *groups.GetAllUseCase,
) *Handler {
	return &Handler{
		presenter:      NewPresenter(bot),
		getAllGroupsUC: getAllGroupsUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "day_type_edit_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "day_type_edit_"), 10, 64)
		h.selectDayTypeExercise(chatID, dayTypeID)
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
