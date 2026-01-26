package exports

import (
	"bytes"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) WriteDoc(chatID int64, buffer *bytes.Buffer) {
	doc := tgbotapi.FileBytes{Name: "workouts.xlsx", Bytes: buffer.Bytes()}
	msg := tgbotapi.NewDocument(chatID, doc)
	p.bot.Send(msg)
}

func (p *Presenter) CannotDoAction(chatID int64, name string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Не удалось выполнить действие '%s' из-за серверной ошибки. Попробуйте позже", name))
	msg.ParseMode = constants.MarkdownParseMode
	p.bot.Send(msg)
}
