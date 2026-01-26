package stats

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/stats"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Handler struct {
	presenter *Presenter

	PeriodStatsUC   *stats.GetPeriodStatsUseCase
	commonPresenter *common.Presenter
}

func NewHandler(bot *tgbotapi.BotAPI, periodStatsUC *stats.GetPeriodStatsUseCase) *Handler {
	return &Handler{
		presenter:       NewPresenter(bot),
		commonPresenter: common.NewPresenter(bot),
		PeriodStatsUC:   periodStatsUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		h.showStatistics(chatID, period)
	}
}

func (h *Handler) RouteMessage(chatID int64, text string) {
	switch {
	case strings.EqualFold(text, "/stats"):
		h.showStatsMenu(chatID)
	}
}

func (h *Handler) showStatistics(chatID int64, period string) {
	res, err := h.PeriodStatsUC.Execute(chatID, period)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.PeriodStatsUC.Name())
		return
	}
	h.presenter.ShowPeriodStats(chatID, res)
}

func (h *Handler) showStatsMenu(chatID int64) {
	h.presenter.ShowStatsMenu(chatID)
}
