package stats

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowPeriodStats(chatID int64, res *dto.PeriodStats) {
	completedWorkouts := res.CompletedWorkouts
	avgTime := res.AvgTime
	cardioTime := res.CardioTime
	isWeek := res.IsWeek
	isMonth := res.IsMonth

	var statsText strings.Builder
	if isWeek {
		statsText.WriteString(messages.StatisticsWeek)
	} else if isMonth {
		statsText.WriteString(messages.StatisticsMonth)
	} else {
		statsText.WriteString(messages.StatisticsAll)
	}
	statsText.WriteString("\n\n")
	statsText.WriteString(messages.EndsWorkouts + fmt.Sprintf(": %d\n", completedWorkouts))
	statsText.WriteString(messages.AvgWorkoutTime + fmt.Sprintf(": %s\n", utils.FormatDuration(avgTime)))
	statsText.WriteString(messages.OverallWorkoutTime + fmt.Sprintf(": %d мин\n", cardioTime))
	text := statsText.String()

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func (p *Presenter) ShowStatsMenu(chatID int64) {
	text := "📊 *Статистика тренировок*\n\n Выберите период:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.StatsWeek, "stats_week"),
			tgbotapi.NewInlineKeyboardButtonData(messages.StatsMonth, "stats_month"),
			tgbotapi.NewInlineKeyboardButtonData(messages.StatsOverall, "stats_all"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}
