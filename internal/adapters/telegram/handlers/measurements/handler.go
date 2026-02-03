package measurements

import (
	measurementsusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/measurements"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter *Presenter

	findAllMeasurementsUC *measurementsusecases.FindAllByUserUseCase
}

func NewHandler(bot *tgbotapi.BotAPI, findAllMeasurementsUC *measurementsusecases.FindAllByUserUseCase) *Handler {
	return &Handler{
		presenter:             NewPresenter(bot),
		findAllMeasurementsUC: findAllMeasurementsUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.EqualFold(data, "measurements_menu"):
		h.showMenu(chatID)
	case strings.HasPrefix(data, "measurements_show_limit_"):
		parts := strings.TrimPrefix(data, "measurements_show_limit_")
		arr := strings.Split(parts, "_")
		limit, _ := strconv.ParseInt(arr[0], 10, 64)
		offset, _ := strconv.ParseInt(arr[1], 10, 64)
		h.showWithLimitAndOffset(chatID, int(limit), int(offset))
	}
}

func (h *Handler) showMenu(chatID int64) {
	h.presenter.showMenu(chatID)
}

func (h *Handler) showWithLimitAndOffset(chatID int64, limit, offset int) {
	res, err := h.findAllMeasurementsUC.Execute(chatID, limit, offset)
	if err != nil {
		return
	}
	h.presenter.showLimitOffset(chatID, limit, offset, res)
}

func (h *Handler) RouteMessage(chatID int64, text string) {
	switch {
	case strings.EqualFold(text, "measurements_menu"):
		h.showMenu(chatID)
	}
}
