package router

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Router) routeMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	fmt.Println("HandleMessage:", text)
	user, _ := r.getUserUC.Execute(chatID)

	switch {
	case text == messages.BackToMenu || text == "/start" || text == "/menu":
		r.sendMainMenu(chatID, message.From)

	case text == messages.StartWorkout || text == "/start_workout":
		r.workoutsHandler.RouteMessage(chatID, "/workouts/start")

	case text == messages.MyWorkouts || text == "/workouts":
		r.workoutsHandler.RouteMessage(chatID, "/workouts")

	case text == messages.Stats || text == "/stats":
		r.statsHandler.RouteMessage(chatID, "/stats")

	case text == messages.Settings || text == "/settings":
		r.settings(chatID)

	case text == messages.HowToUse || text == "/about":
		r.about(chatID)

	case text == messages.Admin || text == "/admin":
		r.admin(chatID, user)

	default:
		r.changesHandler.RouteMessage(chatID, text)
	}
}

func (r *Router) sendMainMenu(chatID int64, from *tgbotapi.User) {
	text := messages.Hello

	user, _ := r.createUserUC.Execute(chatID, from)

	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.StartWorkout),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.MyWorkouts),
		tgbotapi.NewKeyboardButton(messages.Stats),
	))
	rows = append(rows, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(messages.Settings),
		tgbotapi.NewKeyboardButton(messages.HowToUse),
	))

	if user.IsAdmin() {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(messages.Admin),
		))
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	r.bot.Send(msg)
}

func (r *Router) settings(chatID int64) {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.ProgramManagement, "program_management"),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Export, "export_to_excel"),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.Exercises, "exercise_show_all_groups"),
	))
	msg := tgbotapi.NewMessage(chatID, "<b>Выберите действие:</b>")
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	r.bot.Send(msg)
}

func (r *Router) about(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, messages.About)
	msg.ParseMode = constants.HtmlParseMode
	r.bot.Send(msg)
}

func (r *Router) admin(chatID int64, user *models.User) {
	if !user.IsAdmin() {
		return
	}

	msg := tgbotapi.NewMessage(chatID, "<b>👨🏻‍💻 Админ панель</b>")
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Users, "/admin/users"),
		),
	)
	r.bot.Send(msg)
}
