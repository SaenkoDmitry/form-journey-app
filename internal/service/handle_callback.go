package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/templates"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	fmt.Println("HandleCallback:", data)

	switch {
	case strings.HasPrefix(data, "create_workout_"):
		workoutType := strings.TrimPrefix(data, "create_workout_")
		s.createWorkoutDay(chatID, workoutType)

	case strings.HasPrefix(data, "start_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "start_workout_"), 10, 64)
		s.startSpecificWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "start_active_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "start_active_workout_"), 10, 64)
		s.startSpecificWorkout(chatID, workoutID)

	case data == "my_workouts" || data == "create_new_workout":
		s.showMyWorkouts(chatID)

	case strings.HasPrefix(data, "view_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "view_workout_"), 10, 64)
		s.showWorkoutDetails(chatID, workoutID)

	case strings.HasPrefix(data, "confirm_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "confirm_delete_"), 10, 64)
		s.confirmDeleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "delete_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "delete_workout_"), 10, 64)
		s.deleteWorkout(chatID, workoutID)

	case data == "back_to_menu":
		s.sendMainMenu(chatID)

	case strings.HasPrefix(data, "show_progress_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "show_progress_"), 10, 64)
		s.showWorkoutProgress(chatID, workoutID)

	case strings.HasPrefix(data, "continue_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "continue_workout_"), 10, 64)
		s.showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "complete_set_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "complete_set_ex_"), 10, 64)
		s.completeExerciseSet(chatID, exerciseID)

	case strings.HasPrefix(data, "start_timer_"):
		fmt.Println("start_timer_: data: ", data)
		parts := strings.Split(data, "_")
		if len(parts) >= 5 && parts[3] == "ex" {
			seconds, _ := strconv.Atoi(parts[2])
			exerciseID, _ := strconv.ParseInt(parts[4], 10, 64)
			s.startRestTimerWithExercise(chatID, seconds, exerciseID)
		}

	case strings.HasPrefix(data, "prev_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "prev_exercise_"), 10, 64)
		s.moveToPrevExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "next_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "next_exercise_"), 10, 64)
		s.moveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "finish_workout_id_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "finish_workout_id_"), 10, 64)
		s.confirmFinishWorkout(chatID, workoutDayID)

	case strings.HasPrefix(data, "do_finish_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "do_finish_workout_"), 10, 64)
		s.finishWorkoutById(chatID, workoutDayID)

	case strings.HasPrefix(data, "stats_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "stats_workout_"), 10, 64)
		s.showWorkoutStatistics(chatID, workoutID)

	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		s.showStatistics(chatID, period)

	case strings.HasPrefix(data, "change_reps_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_reps_ex_"), 10, 64)
		s.askForNewReps(chatID, exerciseID)

	case strings.HasPrefix(data, "change_weight_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_weight_ex_"), 10, 64)
		s.askForNewWeight(chatID, exerciseID)

	case strings.HasPrefix(data, "change_minutes_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_minutes_ex_"), 10, 64)
		s.askForNewMinutes(chatID, exerciseID)
	}
}

func (s *serviceImpl) showWorkoutProgress(chatID, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Прогресс: %s*\n\n", utils.GetWorkoutNameByID(workoutDay.Name)))

	totalExercises := len(workoutDay.Exercises)
	totalSets := 0
	completedExercises := 0
	completedSets := 0
	for _, exercise := range workoutDay.Exercises {
		totalSets += len(exercise.Sets)
		if exercise.CompletedSets() == len(exercise.Sets) {
			completedExercises++
		}
		completedSets += exercise.CompletedSets()
	}

	text.WriteString(workoutDay.String())

	progressPercent := 0
	if totalSets > 0 {
		progressPercent = (completedSets * 100) / totalSets
	}

	text.WriteString(fmt.Sprintf("\n📈 *Общий прогресс:*\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, totalExercises))
	text.WriteString(fmt.Sprintf("• Подходов: %d/%d\n", completedSets, totalSets))
	text.WriteString(fmt.Sprintf("• Прогресс: %d%%\n", progressPercent))

	barLength := 15
	filled := (progressPercent * barLength) / 100
	progressBar := ""
	for i := 0; i < barLength; i++ {
		if i < filled {
			progressBar += "🏋️‍♂️" // █
		} else {
			progressBar += "░" // ░
		}
	}
	text.WriteString(fmt.Sprintf("• [%s]\n\n", progressBar))

	if workoutDay.EndedAt == nil && completedSets > 0 {
		elapsed := time.Since(workoutDay.StartedAt)
		setsPerMinute := float64(completedSets) / elapsed.Minutes()
		if setsPerMinute > 0 {
			remainingSets := totalSets - completedSets
			remainingMinutes := float64(remainingSets) / setsPerMinute
			text.WriteString(fmt.Sprintf("⏰ *Прогноз окончания:* ~%.0f минут\n", remainingMinutes))
		}
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Продолжить",
				fmt.Sprintf("continue_workout_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) createWorkoutDay(chatID int64, workoutType string) {
	user := s.usersRepo.GetUserByChatID(chatID)

	workoutDay := models.WorkoutDay{
		UserID:    user.ID,
		Name:      workoutType,
		StartedAt: time.Now(),
		Completed: false,
	}

	previousWorkout, _ := s.workoutsRepo.FindPreviousByType(user.ID, workoutType)

	if previousWorkout.ID > 0 {
		fmt.Println("createWorkoutDay: берем настройки количества повторений и веса из последней тренировки")
		for _, exercise := range previousWorkout.Exercises {
			newExercise := models.Exercise{
				Name:          exercise.Name,
				RestInSeconds: exercise.RestInSeconds,
			}
			for _, set := range exercise.Sets {
				newSet := models.Set{
					Reps:    set.GetRealReps(),
					Weight:  set.GetRealWeight(),
					Minutes: set.GetRealMinutes(),
					Index:   set.Index,
				}
				newExercise.Sets = append(newExercise.Sets, newSet)
			}
			workoutDay.Exercises = append(workoutDay.Exercises, newExercise)
		}
	} else {
		fmt.Println("createWorkoutDay: берем настройки количества повторений и веса из preset-ов")
		switch workoutType {
		case constants.LegsAndShouldersWorkoutID:
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetLegExercises()...)
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetShoulderExercises()...)
		case constants.BackAndBicepsWorkoutID:
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetBackExercises()...)
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetBicepsExercises()...)
		case constants.ChestAndTricepsID:
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetChestExercises()...)
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetTricepsExercises()...)
		case constants.CardioID:
			workoutDay.Exercises = append(workoutDay.Exercises, templates.GetCardioExercises()...)
		}
	}

	s.workoutsRepo.Create(&workoutDay)
	s.showCreatedWorkout(chatID, workoutDay.ID)
}

func (s *serviceImpl) showCreatedWorkout(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	var exercisesText strings.Builder
	exercisesText.WriteString(fmt.Sprintf("✅ *Тренировка создана:*\n\n"))
	exercisesText.WriteString(workoutDay.String())
	exercisesText.WriteString("\n Выберите действие:")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Начать тренировку", fmt.Sprintf("start_workout_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_workout_%d", workoutDay.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, exercisesText.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) showWorkoutDetails(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)
	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		s.bot.Send(msg)
		return
	}

	var text strings.Builder
	text.WriteString(workoutDay.String())

	var keyboard tgbotapi.InlineKeyboardMarkup

	if !workoutDay.Completed {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("▶️ Начать тренировку",
					fmt.Sprintf("start_active_workout_%d", workoutID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить",
					fmt.Sprintf("confirm_delete_%d", workoutID)),
			),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить",
					fmt.Sprintf("confirm_delete_%d", workoutID)),
				tgbotapi.NewInlineKeyboardButtonData("Статистика", fmt.Sprintf("stats_workout_%d", workoutID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) confirmDeleteWorkout(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	text := fmt.Sprintf("🗑️ *Удаление тренировки*\n\n"+
		"Вы уверены, что хотите удалить тренировку:\n"+
		"*%s*?\n\n"+
		"❌ Это действие нельзя отменить!", utils.GetWorkoutNameByID(workoutDay.Name))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("delete_workout_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("view_workout_%d", workoutDay.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) deleteWorkout(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	for _, exercise := range workoutDay.Exercises {
		s.setsRepo.Delete(exercise.ID)
	}

	s.exercisesRepo.Delete(workoutID)
	s.workoutsRepo.Delete(&workoutDay)

	msg := tgbotapi.NewMessage(chatID, "✅ Тренировка успешно удалена!")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) startSpecificWorkout(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		s.bot.Send(msg)
		return
	}

	if workoutDay.Completed {
		msg := tgbotapi.NewMessage(chatID, "❌ Эта тренировка уже завершена. Создайте новую или повторите эту.")
		s.bot.Send(msg)
		return
	}

	session := models.WorkoutSession{
		WorkoutDayID:         workoutDay.ID,
		StartedAt:            time.Now(),
		IsActive:             true,
		CurrentExerciseIndex: 0,
	}
	s.sessionsRepo.Create(&session)
	s.showCurrentExerciseSession(chatID, workoutDay.ID)
}

func (s *serviceImpl) showCurrentExerciseSession(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if len(workoutDay.Exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
		s.bot.Send(msg)
		return
	}

	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	exerciseIndex := session.CurrentExerciseIndex
	if exerciseIndex >= len(workoutDay.Exercises) {
		exerciseIndex = 0
	}

	exercise := workoutDay.Exercises[exerciseIndex]

	var text strings.Builder

	exerciseObj, ok := constants.AllExercises[exercise.Name]
	if !ok {
		msg := tgbotapi.NewMessage(chatID, "❌ Упражнение не найдено.")
		s.bot.Send(msg)
		return
	}

	text.WriteString(fmt.Sprintf("<b>%s</b>\n\n", utils.GetWorkoutNameByID(workoutDay.Name)))
	text.WriteString(fmt.Sprintf("<b>Упражнение %d/%d:</b> %s\n\n", exerciseIndex+1, len(workoutDay.Exercises), exerciseObj.GetName()))
	if exerciseObj.GetAccent() != "" {
		text.WriteString(fmt.Sprintf("<b>Акцент:</b> %s\n\n", exerciseObj.GetAccent()))
	}

	for _, set := range exercise.Sets {
		if set.Reps > 0 {
			text.WriteString(fmt.Sprintf("• %s повторов по %s кг: ", set.FormatReps(), set.FormatWeight()))
		}
		if set.Minutes > 0 {
			text.WriteString(fmt.Sprintf("• %s минут: ", set.FormatMinutes()))
		}
		if set.Completed {
			text.WriteString(fmt.Sprintf("✅, %s", set.CompletedAt.Add(3*time.Hour).Format("15:04:05")))
		} else {
			text.WriteString("🚀")
		}
		text.WriteString("\n")
	}

	if hint := exerciseObj.GetHint(); hint != "" {
		text.WriteString(hint)
	}

	changeSettingsButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕➖ Повторения",
			fmt.Sprintf("change_reps_ex_%d", exercise.ID)),
		tgbotapi.NewInlineKeyboardButtonData("⚖️ Вес",
			fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
	)
	if len(exercise.Sets) > 0 && exercise.Sets[0].Minutes > 0 {
		changeSettingsButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⌛ Минуты",
				fmt.Sprintf("change_minutes_ex_%d", exercise.ID)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Завершить подход",
				fmt.Sprintf("complete_set_ex_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏱️ Таймер",
				fmt.Sprintf("start_timer_%d_ex_%d", exercise.RestInSeconds, exercise.ID)),
		),
		changeSettingsButtons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад",
				fmt.Sprintf("prev_exercise_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("➡️ Вперед",
				fmt.Sprintf("next_exercise_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Прогресс",
				fmt.Sprintf("show_progress_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
				fmt.Sprintf("finish_workout_id_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Html"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) completeExerciseSet(chatID int64, exerciseID int64) {
	exercise, _ := s.exercisesRepo.Get(exerciseID)

	nextSet := exercise.NextSet()

	if nextSet.ID != 0 {
		nextSet.Completed = true
		now := time.Now()
		nextSet.CompletedAt = &now
		s.setsRepo.Save(&nextSet)
	}

	exercise, _ = s.exercisesRepo.Get(exerciseID)

	text := fmt.Sprintf("✅ *Подход завершен!*\n\n")
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
}

func (s *serviceImpl) startRestTimerWithExercise(chatID int64, seconds int, exerciseID int64) {
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\n Расслабьтесь и подготовьтесь к следующему подходу!", seconds))

	message, _ := s.bot.Send(msg)

	go func() {
		remaining := seconds

		for remaining > 0 {
			time.Sleep(1 * time.Second)
			remaining--

			if remaining%10 == 0 || remaining <= 5 {
				editMsg := tgbotapi.NewEditMessageText(
					chatID,
					message.MessageID,
					fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\n Расслабьтесь и подготовьтесь к следующему подходу!", remaining),
				)
				s.bot.Send(editMsg)
			}
		}

		editMsg := tgbotapi.NewEditMessageText(
			chatID,
			message.MessageID,
			"🔔 *Время отдыха закончилось!*\n\n Приступайте к следующему подходу! 💪",
		)
		editMsg.ParseMode = "Markdown"

		s.bot.Send(editMsg)

		exercise, _ := s.exercisesRepo.Get(exerciseID)

		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	}()
}

func (s *serviceImpl) moveToExercise(chatID int64, workoutID int64, next bool) {
	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	if session.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Активная сессия не найдена")
		s.bot.Send(msg)
		return
	}

	exercises, _ := s.exercisesRepo.FindAllByWorkoutID(workoutID)

	if len(exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В тренировке нет упражнений")
		s.bot.Send(msg)
		return
	}

	if next {
		session.CurrentExerciseIndex++
	} else {
		session.CurrentExerciseIndex--
	}

	if session.CurrentExerciseIndex < 0 {
		session.CurrentExerciseIndex = 0
		msg := tgbotapi.NewMessage(chatID,
			"Более ранних упражнений в этой тренировке нет")
		s.bot.Send(msg)

		s.showCurrentExerciseSession(chatID, workoutID)
		return
	}

	if session.CurrentExerciseIndex >= len(exercises) {
		session.CurrentExerciseIndex = 0
		msg := tgbotapi.NewMessage(chatID,
			"🎉 Вы завершили все упражнения в этой тренировке!\n\n"+
				"Хотите завершить тренировку или начать заново?")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
					fmt.Sprintf("finish_workout_id_%d", workoutID)),
			),
		)

		msg.ReplyMarkup = keyboard
		s.bot.Send(msg)
		return
	}

	s.sessionsRepo.Save(&session)
	s.showCurrentExerciseSession(chatID, workoutID)
}

func (s *serviceImpl) moveToPrevExercise(chatID int64, workoutID int64) {
	s.moveToExercise(chatID, workoutID, false)
}

func (s *serviceImpl) moveToNextExercise(chatID int64, workoutID int64) {
	s.moveToExercise(chatID, workoutID, true)
}

func (s *serviceImpl) confirmFinishWorkout(chatID int64, workoutDayID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutDayID)

	text := fmt.Sprintf("🏁 *Завершение тренировки*\n\n"+
		"Вы уверены, что хотите завершить тренировку:\n"+
		"*%s*?\n\n"+
		"После завершения вы сможете просмотреть статистику, "+
		"но не сможете добавлять новые подходы.", utils.GetWorkoutNameByID(workoutDay.Name))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, завершить",
				fmt.Sprintf("do_finish_workout_%d", workoutDayID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, продолжить",
				fmt.Sprintf("continue_workout_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) finishWorkoutById(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	now := time.Now()
	workoutDay.Completed = true
	workoutDay.EndedAt = &now
	s.workoutsRepo.Save(&workoutDay)

	s.sessionsRepo.UpdateIsActive(workoutID, false)
	s.showWorkoutStatistics(chatID, workoutID)
}

func (s *serviceImpl) showWorkoutStatistics(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	totalWeight := 0.0
	completedExercises := 0
	totalTime := 0

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Статистика: %s*\n\n", utils.GetWorkoutNameByID(workoutDay.Name)))

	if workoutDay.EndedAt != nil {
		text.WriteString(fmt.Sprintf("⏱️ *Время:* %s\n", utils.BetweenTimes(workoutDay.StartedAt, workoutDay.EndedAt)))
	}

	text.WriteString(fmt.Sprintf("📅 *Дата:* %s\n\n", workoutDay.StartedAt.Add(3*time.Hour).Format("02.01.2006 15:04")))

	for _, exercise := range workoutDay.Exercises {
		if exercise.CompletedSets() == 0 {
			continue
		}

		exerciseObj, ok := constants.AllExercises[exercise.Name]
		if !ok {
			continue
		}

		completedExercises++
		exerciseTime := 0
		exerciseWeight := 0.0
		maxWeight := 0.0

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}
			exerciseWeight += float64(set.GetRealWeight()) * float64(set.GetRealReps())
			exerciseTime += set.GetRealMinutes()
			maxWeight = max(maxWeight, float64(set.GetRealWeight()))
		}
		totalWeight += exerciseWeight
		totalTime += exerciseTime

		lastSet := exercise.Sets[len(exercise.Sets)-1]
		text.WriteString(fmt.Sprintf("• *%s:* \n", exerciseObj.GetName()))
		if lastSet.GetRealReps() > 0 {
			text.WriteString(fmt.Sprintf("  • Выполнено: %d из %d повторений\n", exercise.CompletedSets(), len(exercise.Sets)))
			text.WriteString(fmt.Sprintf("  • Рабочий вес: %d \\* %.0f кг \n", lastSet.GetRealReps(), lastSet.GetRealWeight()))
			text.WriteString(fmt.Sprintf("  • Общий вес: %.0f кг \n\n", exerciseWeight))
		} else if lastSet.GetRealMinutes() > 0 {
			text.WriteString(fmt.Sprintf("  • Общее время: %d минут \n\n", exerciseTime))
		}
	}

	text.WriteString(fmt.Sprintf("📈 *Итого:*\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, len(workoutDay.Exercises)))
	if totalWeight > 0 {
		text.WriteString(fmt.Sprintf("• Общий тоннаж: %.0f кг\n", totalWeight))
	}
	if totalTime > 0 {
		text.WriteString(fmt.Sprintf("• Общее время: %d минут\n", totalTime))
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)
}

func (s *serviceImpl) showStatistics(chatID int64, period string) {
	user := s.usersRepo.GetUserByChatID(chatID)

	workouts, _ := s.workoutsRepo.Find(user.ID)

	completedStrengthWorkouts := 0
	sumStrengthTime := time.Duration(0)
	cardioTime := 0
	for _, w := range workouts {
		if !w.Completed {
			continue
		}
		switch period {
		case "week":
			if time.Since(w.StartedAt).Hours() > 7*24 {
				continue
			}
		case "month":
			if time.Since(w.StartedAt).Hours() > 30*24 {
				continue
			}
		default:
		}

		// not cardio
		if w.Name != constants.CardioID {
			fmt.Println("not cardio count exercises:", len(w.Exercises))
			completedStrengthWorkouts++
			sumStrengthTime += w.EndedAt.Sub(*&w.StartedAt)
		} else {
			// cardio
			fmt.Println("cardio count exercises:", len(w.Exercises))
			for _, e := range w.Exercises {
				if len(e.Sets) == 0 {
					continue
				}
				for _, s := range e.Sets {
					if !s.Completed {
						continue
					}
					cardioTime += s.GetRealMinutes()
				}
			}
		}
	}
	avgTime := sumStrengthTime / time.Duration(completedStrengthWorkouts)

	var statsText strings.Builder
	statsText.WriteString("📅 *Статистика за неделю*\n\n")
	statsText.WriteString(fmt.Sprintf("✅ Силовых тренировок: %d\n", completedStrengthWorkouts))
	statsText.WriteString(fmt.Sprintf("⏱️ Среднее время сил. тренировки: %s\n", utils.FormatDuration(avgTime)))
	statsText.WriteString(fmt.Sprintf("🫀 Общее время кардио: %d мин\n", cardioTime))

	msg := tgbotapi.NewMessage(chatID, statsText.String())
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewReps(chatID int64, exerciseID int64) {
	s.userStates[chatID] = fmt.Sprintf("awaiting_reps_%d", exerciseID)
	msg := tgbotapi.NewMessage(chatID, "➕➖ Введите новое число повторений:")
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewWeight(chatID int64, exerciseID int64) {
	s.userStates[chatID] = fmt.Sprintf("awaiting_weight_%d", exerciseID)
	msg := tgbotapi.NewMessage(chatID, "⚖️ Введите новый вес (в кг):")
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewMinutes(chatID int64, exerciseID int64) {
	s.userStates[chatID] = fmt.Sprintf("awaiting_minutes_%d", exerciseID)
	msg := tgbotapi.NewMessage(chatID, "⌛ Введите новое время (мин):")
	s.bot.Send(msg)
}
