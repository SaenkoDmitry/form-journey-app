package daytypes

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowSelectDayTypeDialog(chatID int64, dayTypeID int64, res *dto.ExerciseGroupTypeList) {
	groups := res.Groups

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(group.Name, fmt.Sprintf("exercise_select_for_program_day_%d_%s", dayTypeID, group.Code)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, messages.SelectGroupOfMuscle)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}
