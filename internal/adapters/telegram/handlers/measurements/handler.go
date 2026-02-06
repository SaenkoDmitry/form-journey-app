package measurements

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	measurementsusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/measurements"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	findAllMeasurementsUC   *measurementsusecases.FindAllByUserUseCase
	getMeasurementByIDUC    *measurementsusecases.GetByIDUseCase
	deleteMeasurementByIDUC *measurementsusecases.DeleteByIDUseCase
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	findAllMeasurementsUC *measurementsusecases.FindAllByUserUseCase,
	getMeasurementByIDUC *measurementsusecases.GetByIDUseCase,
	deleteMeasurementByIDUC *measurementsusecases.DeleteByIDUseCase,
) *Handler {
	return &Handler{
		presenter:               NewPresenter(bot),
		commonPresenter:         common.NewPresenter(bot),
		findAllMeasurementsUC:   findAllMeasurementsUC,
		getMeasurementByIDUC:    getMeasurementByIDUC,
		deleteMeasurementByIDUC: deleteMeasurementByIDUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.EqualFold(data, "measurements_menu"):
		h.showMenu(chatID)
	case strings.HasPrefix(data, "measurements_view_"):
		measurementID, _ := strconv.ParseInt(strings.TrimPrefix(data, "measurements_view_"), 10, 64)
		h.viewMeasurement(chatID, measurementID)
	case strings.HasPrefix(data, "measurements_delete_"):
		measurementID, _ := strconv.ParseInt(strings.TrimPrefix(data, "measurements_delete_"), 10, 64)
		h.deleteMeasurement(chatID, measurementID)
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

func (h *Handler) viewMeasurement(chatID int64, measurementID int64) {
	res, err := h.getMeasurementByIDUC.Execute(measurementID)
	if err != nil {
		return
	}
	h.presenter.viewMeasurement(chatID, res)
}

func (h *Handler) deleteMeasurement(chatID int64, measurementID int64) {
	err := h.deleteMeasurementByIDUC.Execute(measurementID)
	if err != nil {
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.SuccessfullyDeletedMeasurement)
	h.showMenu(chatID)
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
