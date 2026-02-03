package measurements

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p Presenter) showMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"<b>Выберите действие:</b>")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новое", "change_add_new_measurement"),
			formatMoveToButton("📋 История", 3, 0),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить последнее добавленное", "measurements_delete_last"),
		//),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Export, "export_measurements_to_excel"),
		),
	)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p Presenter) showLimitOffset(chatID int64, limit, offset int, result *dto.FindWithOffsetLimitMeasurement) {
	measurementObjs := result.Items
	count := result.Count

	shoulders := make([]string, 0, len(measurementObjs))
	chests := make([]string, 0, len(measurementObjs))
	handLeft := make([]string, 0, len(measurementObjs))
	handRight := make([]string, 0, len(measurementObjs))
	waists := make([]string, 0, len(measurementObjs))
	buttocks := make([]string, 0, len(measurementObjs))
	hipLeft := make([]string, 0, len(measurementObjs))
	hipRight := make([]string, 0, len(measurementObjs))
	calfLeft := make([]string, 0, len(measurementObjs))
	calfRight := make([]string, 0, len(measurementObjs))
	weights := make([]string, 0, len(measurementObjs))

	var from, to string
	if len(measurementObjs) > 0 {
		from = measurementObjs[len(measurementObjs)-1].CreatedAt
		to = measurementObjs[0].CreatedAt
	}
	for i := len(measurementObjs) - 1; i >= 0; i-- {
		m := measurementObjs[i]
		shoulders = append(shoulders, m.Shoulders)
		chests = append(chests, m.Chest)
		handLeft = append(handLeft, m.HandLeft)
		handRight = append(handRight, m.HandRight)
		waists = append(waists, m.Waist)
		buttocks = append(buttocks, m.Buttocks)
		hipLeft = append(hipLeft, m.HipLeft)
		hipRight = append(hipRight, m.HipRight)
		calfLeft = append(calfLeft, m.CalfLeft)
		calfRight = append(calfRight, m.CalfRight)
		weights = append(weights, m.Weight)
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"<b>%s за период \n"+
			"📆 %s – %s</b>\n\n"+
			"• <u>Плечи (см)</u>: %s\n\n"+
			"• <u>Грудь (см)</u>: %s\n\n"+
			"• <u>Рука левая (см)</u>: %s\n\n"+
			"• <u>Рука правая (см)</u>: %s\n\n"+
			"• <u>Талия (см)</u>: %s\n\n"+
			"• <u>Ягодицы (см)</u>: %s\n\n"+
			"• <u>Бедро левое (см)</u>: %s\n\n"+
			"• <u>Бедро правое (см)</u>: %s\n\n"+
			"• <u>Икра левая (см)</u>: %s\n\n"+
			"• <u>Икра правая (см)</u>: %s\n\n"+
			"• <u>Вес (кг)</u>: %s",
		messages.Measurements,
		from, to,
		strings.Join(shoulders, delimiter),
		strings.Join(chests, delimiter),
		strings.Join(handLeft, delimiter),
		strings.Join(handRight, delimiter),
		strings.Join(waists, delimiter),
		strings.Join(buttocks, delimiter),
		strings.Join(hipLeft, delimiter),
		strings.Join(hipRight, delimiter),
		strings.Join(calfLeft, delimiter),
		strings.Join(calfRight, delimiter),
		strings.Join(weights, delimiter),
	))
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})
	if offset+limit < count {
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1], formatMoveToButton(messages.Earlier, limit, offset+limit))
	}
	if offset-limit >= 0 {
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1], formatMoveToButton(messages.Later, limit, offset-limit))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, "measurements_menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func formatMoveToButton(text string, limit, offset int) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, fmt.Sprintf("measurements_show_limit_%d_%d", limit, offset))
}

const (
	delimiter = " -> "
)
