package exports

import (
	exportusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exports"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Handler struct {
	presenter       *Presenter
	exportToExcelUC *exportusecases.ExportToExcelUseCase
}

func NewHandler(bot *tgbotapi.BotAPI, exportToExcelUC *exportusecases.ExportToExcelUseCase) *Handler {
	return &Handler{
		presenter:       NewPresenter(bot),
		exportToExcelUC: exportToExcelUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "export_to_excel"):
		h.exportToExcel(chatID)
	}
}

func (h *Handler) exportToExcel(chatID int64) {
	buffer, err := h.exportToExcelUC.Execute(chatID)
	if err != nil {
		h.presenter.CannotDoAction(chatID, h.exportToExcelUC.Name())
		return
	}
	h.presenter.WriteDoc(chatID, buffer)
}
