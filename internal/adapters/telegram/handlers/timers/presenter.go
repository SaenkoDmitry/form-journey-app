package timers

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot             *tgbotapi.BotAPI
	commonPresenter *common.Presenter
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot, commonPresenter: common.NewPresenter(bot)}
}

func (p *Presenter) ShowTimerIsNotSupported(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, messages.RestNotSupported)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func (p *Presenter) ShowCreatedTimer(chatID int64, seconds int, res *dto.StartTimer, doWhenTimerExpired func()) {
	newTimerID := res.NewTimerID
	remainingCh := res.RemainingCh

	var message tgbotapi.Message

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.RestTimer, seconds))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.CancelTimer, fmt.Sprintf("timer_unpin_and_cancel_%s", newTimerID)),
	))
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	message, _ = p.bot.Send(msg)
	p.commonPresenter.PinMessage(chatID, message)

	go func() {
		for remaining := range remainingCh {
			if remaining%10 == 0 || remaining <= 20 {
				p.handleTimerTick(chatID, remaining, message, keyboard)
			}
		}

		if res.IsStopped(chatID, newTimerID) {
			p.handleTimerStopped(chatID, message)
			return
		}

		p.handleTimerExpired(chatID, newTimerID, message, res.StopTimer)
		doWhenTimerExpired()
	}()
}

func (p *Presenter) handleTimerExpired(chatID int64, timerID string, message tgbotapi.Message, stopTimer func(chatID int64, timerID string)) {
	fmt.Println("expired")
	stopTimer(chatID, timerID)
	editMsg := tgbotapi.NewEditMessageText(
		chatID,
		message.MessageID,
		messages.RestIsEnded,
	)
	editMsg.ParseMode = constants.HtmlParseMode
	editMessage, _ := p.bot.Send(editMsg)
	p.commonPresenter.UnpinMessage(chatID, editMessage)
}

func (p *Presenter) handleTimerStopped(chatID int64, message tgbotapi.Message) {
	fmt.Println("stopped")
	p.commonPresenter.UnpinMessage(chatID, message)
	editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, messages.TimerCanceled)
	editMsg.ParseMode = constants.HtmlParseMode
	p.bot.Send(editMsg)
}

func (p *Presenter) handleTimerTick(chatID int64, remaining int, message tgbotapi.Message, keyboard tgbotapi.InlineKeyboardMarkup) {
	editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, message.MessageID,
		fmt.Sprintf(messages.RestTimer, remaining), keyboard)
	editMsg.ParseMode = constants.HtmlParseMode
	message, _ = p.bot.Send(editMsg)
}
