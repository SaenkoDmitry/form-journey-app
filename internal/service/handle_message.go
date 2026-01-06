package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	fmt.Println("HandleMessage:", text)

	user, _ := s.usersRepo.GetUser(chatID, message.From.UserName)

	switch {
	case text == "/start" || text == "/menu" || text == "🔙 В меню":
		s.sendMainMenu(chatID)

	case text == "/start_workout" || text == "▶️ Начать тренировку":
		s.showWorkoutTypeMenu(chatID)

	case text == "/stats" || text == "📊 Статистика":
		s.showStatsMenu(chatID, user.ID)

	case text == "📋 Мои тренировки" || text == "/workouts":
		s.showMyWorkouts(chatID)

	default:
		s.handleState(chatID, text)
	}
}

func (s *serviceImpl) sendMainMenu(chatID int64) {
	text := "🏋️‍♂️ *Добро пожаловать в Бот для тренировок!* \n\n Выберите действие:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("▶️ Начать тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Мои тренировки"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
	)
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showWorkoutTypeMenu(chatID int64) {
	text := "Выберите тип тренировки:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(constants.LegsAndShouldersWorkoutName, "create_workout_legs_and_shoulders"),
			tgbotapi.NewInlineKeyboardButtonData(constants.BackAndBicepsWorkoutName, "create_workout_back_and_biceps"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(constants.ChestAndTricepsName, "create_workout_chest_and_triceps"),
			tgbotapi.NewInlineKeyboardButtonData(constants.CardioName, "create_workout_cardio"),
		),
		// tgbotapi.NewInlineKeyboardRow(
		//  tgbotapi.NewInlineKeyboardButtonData("💪 Руки", "create_workout_arms"),
		// 	tgbotapi.NewInlineKeyboardButtonData("🌀 Плечи", "create_workout_shoulders"),
		// 	tgbotapi.NewInlineKeyboardButtonData("🫀 Кардио", "create_workout_cardio"),
		// ),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showMyWorkouts(chatID int64) {
	user := s.usersRepo.GetUserByChatID(chatID)

	workouts, _ := s.workoutsRepo.Find(user.ID)

	if len(workouts) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 У вас пока нет созданных тренировок.\n\nСоздайте первую тренировку!")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
			),
		)
		msg.ReplyMarkup = keyboard
		s.bot.Send(msg)
		return
	}

	text := "📋 *Ваши тренировки:*\n\n"
	for i, workout := range workouts {
		status := "🟡 Активна"
		if workout.Completed {
			status = "✅ Завершена"
		}
		date := workout.StartedAt.Format("02.01.2006")
		text += fmt.Sprintf("%d. *%s* - %s\n   📅 %s\n\n",
			i+1, utils.GetWorkoutNameByID(workout.Name), status, date)
	}

	text += "Выберите тренировку для просмотра:"

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, workout := range workouts {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rowIndex := len(rows) - 1
		buttonText := fmt.Sprintf("%s %d", utils.GetWorkoutNameByID(workout.Name), i+1)
		rows[rowIndex] = append(rows[rowIndex],
			tgbotapi.NewInlineKeyboardButtonData(buttonText,
				fmt.Sprintf("view_workout_%d", workout.ID)))
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
	})

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showStatsMenu(chatID int64, userID int64) {
	text := "📊 *Статистика тренировок*\n\n Выберите период:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 За неделю", "stats_week"),
			tgbotapi.NewInlineKeyboardButtonData("🗓️ За месяц", "stats_month"),
			tgbotapi.NewInlineKeyboardButtonData("📈 Общая", "stats_all"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) handleState(chatID int64, text string) {
	state, exists := s.userStates[chatID]
	if !exists {
		return
	}

	switch {
	case strings.HasPrefix(state, "awaiting_reps_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 3 {
			exerciseID, _ := strconv.ParseInt(parts[2], 10, 64)

			reps, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат числа повторений. Введите целое число (например: 42)")
				s.bot.Send(msg)
				return
			}

			exercise, _ := s.exercisesRepo.Get(exerciseID)

			nextSet := exercise.NextSet()
			if nextSet.ID != 0 {
				nextSet.FactReps = int(reps)
				if int(reps) != nextSet.Reps {
					nextSet.FactReps = int(reps)
				} else {
					nextSet.FactReps = 0
				}
				s.setsRepo.Save(&nextSet)

				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
					"✅ Количество повторений обновлено: %d раз(а) для подхода №%d",
					reps, nextSet.Index,
				))
				s.bot.Send(msg)
			}

			s.userStates[chatID] = ""

			s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
		}
	case strings.HasPrefix(state, "awaiting_weight_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 3 {
			exerciseID, _ := strconv.ParseInt(parts[2], 10, 64)

			weight, err := strconv.ParseFloat(text, 32)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат веса. Введите число (например: 42.5)")
				s.bot.Send(msg)
				return
			}

			exercise, _ := s.exercisesRepo.Get(exerciseID)

			nextSet := exercise.NextSet()
			if nextSet.ID != 0 {
				if float32(weight) != nextSet.Weight {
					nextSet.FactWeight = float32(weight)
				} else {
					nextSet.FactWeight = float32(0)
				}
				s.setsRepo.Save(&nextSet)

				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
					"✅ Вес обновлен: %.1f кг для подхода №%d",
					weight, nextSet.Index,
				))
				s.bot.Send(msg)
			}

			s.userStates[chatID] = ""

			s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
		}

	case strings.HasPrefix(state, "awaiting_minutes_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 3 {
			exerciseID, _ := strconv.ParseInt(parts[2], 10, 64)

			minutes, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат минут. Введите число (например: 42)")
				s.bot.Send(msg)
				return
			}

			exercise, _ := s.exercisesRepo.Get(exerciseID)

			nextSet := exercise.NextSet()
			if nextSet.ID != 0 {
				if int(minutes) != nextSet.Minutes {
					nextSet.FactMinutes = int(minutes)
				} else {
					nextSet.FactMinutes = int(0)
				}
				s.setsRepo.Save(&nextSet)

				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
					"✅ Время обновлено: %d минут для подхода №%d",
					minutes, nextSet.Index,
				))
				s.bot.Send(msg)
			}

			s.userStates[chatID] = ""

			s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
		}
	}
}
