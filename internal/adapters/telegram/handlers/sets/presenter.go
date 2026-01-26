package sets

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func (p *Presenter) ShowOneMore(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, messages.SetAdded)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func (p *Presenter) ShowYouCannotDeleteSet(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}
