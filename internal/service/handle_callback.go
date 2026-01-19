package service

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	fmt.Println("HandleCallback:", data)

	switch {
	case data == "back_to_menu":
		s.sendMainMenu(chatID, callback.From)

	case data == "/admin/users":
		s.users(chatID, user)

	case strings.HasPrefix(data, "program_"):
		s.programCases(data, chatID, user.ID)

	case strings.HasPrefix(data, "workout_"):
		s.workoutCases(data, chatID, user.ID)

	case strings.HasPrefix(data, "timer_"):
		s.timerCases(data, chatID)

	case strings.HasPrefix(data, "set_"):
		s.setCases(data, chatID)

	case strings.HasPrefix(data, "exercise_"):
		s.exerciseCases(data, chatID)

	case strings.HasPrefix(data, "change_"):
		s.changeCases(data, chatID)

	case strings.HasPrefix(data, "stats_"):
		s.statsCases(data, chatID, user.ID)
	}
}

func (s *serviceImpl) selectExerciseForCurrentWorkout(chatID int64, workoutID int64, exerciseGroupCode string) {
	method := "selectExerciseForCurrentWorkout"
	group, err := s.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return
	}

	text := fmt.Sprintf("<b>Тип:</b> %s \n\n <b>Выберите упражнение из списка:</b>", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	exerciseTypes, err := s.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return
	}

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("exercise_add_specific_for_current_workout_%d_%d", workoutID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)

	msg := tghelpers.NewMessageBuilder().
		WithChatID(chatID).
		WithText(text).
		WithReplyMarkup(rows).
		Build()
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) GetUserByChatID(chatID int64) (*models.User, error) {
	method := "GetUserByChatID"
	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, users.NotFoundUserErr) {
			msg := tghelpers.NewMessageBuilder().
				WithChatID(chatID).
				WithText("Сначала создайте пользователя в боте, через команду /start").
				Build()
			_, _ = tghelpers.SendMessage(s.bot, msg, method)
		}
	}
	return user, nil
}
