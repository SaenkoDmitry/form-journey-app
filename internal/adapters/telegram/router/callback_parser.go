package router

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (r *Router) routeCallback(callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	fmt.Println("HandleCallback:", data)

	switch {
	case data == "back_to_menu":
		r.sendMainMenu(chatID, callbackQuery.From)

	case data == "/admin/users":
		// todo s.users(chatID, user)

	case strings.HasPrefix(data, "program_"):
		r.programsHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "workout_"):
		r.workoutsHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "timer_"):
		r.timersHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "set_"):
		r.setsHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "exercise_"):
		r.exercisesHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "change_"):
		r.changesHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "stats_"):
		r.statsHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "export_"):
		r.exportsHandler.RouteCallback(chatID, data)

	case strings.HasPrefix(data, "day_type_"):
		r.dayTypesHandler.RouteCallback(chatID, data)
	}
}
