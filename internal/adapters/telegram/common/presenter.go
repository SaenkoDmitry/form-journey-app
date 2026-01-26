package common

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) PinMessage(chatID int64, message tgbotapi.Message) {
	pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	p.bot.Request(pinChatMessageConfig)
}

func (p *Presenter) UnpinMessage(chatID int64, message tgbotapi.Message) {
	unpinChatMessageConfig := tgbotapi.UnpinChatMessageConfig{
		ChatID:    chatID,
		MessageID: message.MessageID,
	}
	p.bot.Request(unpinChatMessageConfig)
}

func (p *Presenter) HandleInternalError(err error, chatID int64, name string) {
	if errors.Is(err, users.NotFoundUserErr) {
		msg := tgbotapi.NewMessage(chatID, messages.FirstCreateUser)
		msg.ParseMode = constants.MarkdownParseMode
		p.bot.Send(msg)
	}

	text := fmt.Sprintf("❌ Не удалось выполнить действие '%s' "+
		"из-за серверной ошибки. Попробуйте позже", name)
	fmt.Println(text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	p.bot.Send(msg)
}

func (p *Presenter) SendSimpleHtmlMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}
