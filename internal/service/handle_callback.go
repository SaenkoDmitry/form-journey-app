package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/templates"
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

	}
}

func (s *serviceImpl) showWorkoutProgress(chatID, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Прогресс тренировки: %s*\n\n", workoutDay.Name))

	totalExercises := len(workoutDay.Exercises)
	completedExercises := 0
	totalSets := 0
	completedSets := 0

	text.WriteString(workoutDay.String())

	progressPercent := 0
	if totalSets > 0 {
		progressPercent = (completedSets * 100) / totalSets
	}

	text.WriteString(fmt.Sprintf("\n📈 *Общий прогресс:*\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, totalExercises))
	text.WriteString(fmt.Sprintf("• Подходов: %d/%d\n", completedSets, totalSets))
	text.WriteString(fmt.Sprintf("• Прогресс: %d%%\n", progressPercent))

	barLength := 10
	filled := (progressPercent * barLength) / 100
	progressBar := ""
	for i := 0; i < barLength; i++ {
		if i < filled {
			progressBar += "█"
		} else {
			progressBar += "░"
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
			// tgbotapi.NewInlineKeyboardButtonData("📊 Детали",
			// 	fmt.Sprintf("detailed_stats_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 К тренировке",
				fmt.Sprintf("view_workout_%d", workoutID)),
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
	switch workoutType {
	case "legs":
		workoutDay.Exercises = templates.GetLegExercises()
	case "back":
		workoutDay.Exercises = templates.GetBackExercises()
	case "chest_and_shoulders":
		workoutDay.Exercises = append(workoutDay.Exercises, templates.GetChestExercises()...)
		workoutDay.Exercises = append(workoutDay.Exercises, templates.GetShoulderExercises()...)
		// case "shoulders":
		// 	workoutDay.Exercises = getShoulderExercises()
		// case "cardio":
		// 	workoutDay.Exercises = getCardioExercises()
	}

	s.workoutsRepo.Create(&workoutDay)
	s.showCreatedWorkout(chatID, workoutDay.ID)
}

func (s *serviceImpl) showCreatedWorkout(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	var exercisesText strings.Builder
	exercisesText.WriteString(fmt.Sprintf("✅ *Тренировка создана: %s*\n\n", workoutDay.Name))
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
		"❌ Это действие нельзя отменить!", workoutDay.Name)

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

	text.WriteString(fmt.Sprintf("🏋️‍♂️ *Тренировка: %s*\n\n", workoutDay.Name))
	text.WriteString(fmt.Sprintf("*Упражнение %d/%d:* %s\n\n", exerciseIndex+1, len(workoutDay.Exercises), exercise.Name))
	text.WriteString(fmt.Sprintf("Выполнено: %d из %d подходов\n\n", exercise.CompletedSets(), len(exercise.Sets)))
	for _, set := range exercise.Sets {
		text.WriteString(fmt.Sprintf("%d повторов по %.0f кг: ", set.Reps, set.Weight))
		if set.Completed {
			text.WriteString(fmt.Sprintf("✅, %s", set.CompletedAt.Add(3*time.Hour).Format("15:04:05")))
		} else {
			text.WriteString("🚀")
		}
		text.WriteString("\n")
	}
	text.WriteString("\n\n *Что делаем?*")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Завершить подход",
				fmt.Sprintf("complete_set_ex_%d", exercise.ID)),
		),
		// tgbotapi.NewInlineKeyboardRow(
		// 	tgbotapi.NewInlineKeyboardButtonData("➕ Больше повторений",
		// 		fmt.Sprintf("add_reps_ex_%d", exercise.ID)),
		// 	tgbotapi.NewInlineKeyboardButtonData("⚖️ Изменить вес",
		// 		fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
		// ),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏸️ Таймер отдыха",
				fmt.Sprintf("start_timer_%d_ex_%d", exercise.RestInSeconds, exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("➡️ След. упр-е",
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
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) completeExerciseSet(chatID int64, exerciseID int64) {
	exercise, _ := s.exercisesRepo.Get(exerciseID)

	nextSet := exercise.Next()

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

func (s *serviceImpl) moveToNextExercise(chatID int64, workoutID int64) {
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

	session.CurrentExerciseIndex++

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
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Начать заново",
					fmt.Sprintf("restart_workout_%d", workoutID)),
				tgbotapi.NewInlineKeyboardButtonData("🔙 К первому",
					fmt.Sprintf("first_exercise_%d", workoutID)),
			),
		)

		msg.ReplyMarkup = keyboard
		s.bot.Send(msg)
		return
	}

	s.sessionsRepo.Save(&session)
	s.showCurrentExerciseSession(chatID, workoutID)
}

func (s *serviceImpl) confirmFinishWorkout(chatID int64, workoutDayID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutDayID)

	text := fmt.Sprintf("🏁 *Завершение тренировки*\n\n"+
		"Вы уверены, что хотите завершить тренировку:\n"+
		"*%s*?\n\n"+
		"После завершения вы сможете просмотреть статистику, "+
		"но не сможете добавлять новые подходы.", workoutDay.Name)

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

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Статистика: %s*\n\n", workoutDay.Name))

	if workoutDay.EndedAt != nil {
		duration := workoutDay.EndedAt.Sub(workoutDay.StartedAt)
		text.WriteString(fmt.Sprintf("⏱️ *Время:* %s\n", formatDuration(duration)))
	}

	text.WriteString(fmt.Sprintf("📅 *Дата:* %s\n\n", workoutDay.StartedAt.Add(3*time.Hour).Format("02.01.2006 15:04")))

	for _, exercise := range workoutDay.Exercises {
		if exercise.CompletedSets() == 0 {
			continue
		}

		completedExercises++
		exerciseWeight := 0.0
		maxWeight := 0.0

		for _, set := range exercise.Sets {
			if !set.Completed {
				continue
			}
			exerciseWeight += float64(set.Weight) * float64(set.Reps)
			maxWeight = max(maxWeight, float64(set.Weight))
		}
		totalWeight += exerciseWeight

		text.WriteString(fmt.Sprintf("• *%s:* %d из %d повторений (макс вес %.0f кг, общий вес %.0f кг)\n\n",
			exercise.Name, exercise.CompletedSets(), len(exercise.Sets), maxWeight, exerciseWeight))
	}

	text.WriteString(fmt.Sprintf("📈 *Итого:*\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, len(workoutDay.Exercises)))
	text.WriteString(fmt.Sprintf("• Общий тоннаж: %.0f кг\n", totalWeight))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад",
				fmt.Sprintf("view_workout_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dч %dмин", hours, minutes)
	}
	return fmt.Sprintf("%dмин %dсек", minutes, seconds)
}

func (s *serviceImpl) showStatistics(chatID int64, period string) {
	var statsText string
	user := s.usersRepo.GetUserByChatID(chatID)

	workouts, _ := s.workoutsRepo.Find(user.ID)

	completedWorkouts := 0
	sumTime := 0
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
		completedWorkouts++
		sumTime += int((w.EndedAt.Sub(*&w.StartedAt)).Minutes())
	}
	avgTime := float64(sumTime) / float64(completedWorkouts)

	statsText = fmt.Sprintf("📅 *Статистика за неделю*\n\n✅ Тренировок: %d\n⏱️ Среднее время: %.0f мин", completedWorkouts, avgTime)

	msg := tgbotapi.NewMessage(chatID, statsText)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)
}
