package tghelpers

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageBuilder struct {
	ChatID      int64
	Text        string
	ParseMode   string
	ReplyMarkup *tgbotapi.InlineKeyboardMarkup
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

func (m *MessageBuilder) WithChatID(chatID int64) *MessageBuilder {
	m.ChatID = chatID
	return m
}

func (m *MessageBuilder) WithText(text string) *MessageBuilder {
	m.Text = text
	return m
}

func (m *MessageBuilder) WithParseMode(mode string) *MessageBuilder {
	m.ParseMode = mode
	return m
}

func (m *MessageBuilder) WithReplyMarkup(rows [][]tgbotapi.InlineKeyboardButton) *MessageBuilder {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	m.ReplyMarkup = &keyboard
	return m
}

func (m *MessageBuilder) Build() tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(m.ChatID, m.Text)
	if m.ParseMode == "" { // default
		msg.ParseMode = constants.HtmlParseMode
	}
	if m.ReplyMarkup != nil {
		msg.ReplyMarkup = m.ReplyMarkup
	}
	return msg
}
