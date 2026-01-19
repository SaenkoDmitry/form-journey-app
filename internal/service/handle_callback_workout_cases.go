package service

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/service/tghelpers"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

func (s *serviceImpl) workoutCases(data string, chatID, userID int64) {
	switch {
	case strings.HasPrefix(data, "workout_create_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_create_"), 10, 64)
		s.createWorkoutDay(chatID, userID, dayTypeID)

	case strings.HasPrefix(data, "workout_start_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_start_"), 10, 64)
		s.startSpecificWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "workout_show_my"):
		if data == "workout_show_my" {
			s.showMyWorkouts(chatID, 0)
			return
		}
		offset, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_my_"), 10, 64)
		s.showMyWorkouts(chatID, int(offset))

	case strings.HasPrefix(data, "workout_confirm_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_confirm_delete_"), 10, 64)
		s.confirmDeleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "workout_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_delete_"), 10, 64)
		s.deleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "workout_continue_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_continue_"), 10, 64)
		s.showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "workout_confirm_finish_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_confirm_finish_"), 10, 64)
		s.confirmFinishWorkout(chatID, workoutDayID)

	case strings.HasPrefix(data, "workout_finish_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_finish_"), 10, 64)
		s.finishWorkoutById(chatID, workoutDayID)

	case strings.HasPrefix(data, "workout_stats_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_stats_"), 10, 64)
		s.showWorkoutStatistics(chatID, workoutID)

	case strings.HasPrefix(data, "workout_show_progress_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_progress_"), 10, 64)
		s.showWorkoutProgress(chatID, workoutID)

	case strings.HasPrefix(data, "workout_show_by_user_id_"):
		userID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_by_user_id_"), 10, 64)
		s.showWorkoutsByUser(chatID, userID)
	}
}

func (s *serviceImpl) showWorkoutProgress(chatID, workoutID int64) {
	method := "showWorkoutProgress"
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	// calc stats
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
	progressPercent := 0
	if totalSets > 0 {
		progressPercent = (completedSets * 100) / totalSets
	}
	//

	var text strings.Builder
	text.WriteString(workoutDay.String())
	text.WriteString(fmt.Sprintf("\n📈 <b>Общий прогресс:</b>\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, totalExercises))
	text.WriteString(fmt.Sprintf("• Подходов: %d/%d\n", completedSets, totalSets))
	text.WriteString(fmt.Sprintf("• Прогресс: %d%%\n", progressPercent))

	barLength := 13
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
			text.WriteString(fmt.Sprintf("⏰ <b>Прогноз окончания:</b> ~%.0f минут\n", remainingMinutes))
		}
	}

	oneMoreExerciseButton := tgbotapi.NewInlineKeyboardButtonData("➕ Еще упражнение", fmt.Sprintf("exercise_add_for_current_workout_%d", workoutID))
	deleteExerciseButton := tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("workout_confirm_delete_%d", workoutID))
	showMyWorkoutsButton := tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "workout_show_my")
	statsWorkoutButton := tgbotapi.NewInlineKeyboardButtonData(messages.Stats, fmt.Sprintf("workout_stats_%d", workoutID))

	var toWorkoutButton tgbotapi.InlineKeyboardButton
	if _, err := s.sessionsRepo.GetByWorkoutID(workoutID); err != nil {
		toWorkoutButton = tgbotapi.NewInlineKeyboardButtonData("▶️ Начать", fmt.Sprintf("workout_start_%d", workoutDay.ID))
	} else {
		toWorkoutButton = tgbotapi.NewInlineKeyboardButtonData("▶️ К тренировке", fmt.Sprintf("exercise_show_current_%d", workoutID))
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	if !workoutDay.Completed {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(oneMoreExerciseButton, deleteExerciseButton),
			tgbotapi.NewInlineKeyboardRow(toWorkoutButton),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(statsWorkoutButton),
			tgbotapi.NewInlineKeyboardRow(showMyWorkoutsButton, deleteExerciseButton),
		)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) createWorkoutDay(chatID, userID int64, dayTypeID int64) {
	method := "createWorkoutDay"

	workoutDay := models.WorkoutDay{
		UserID:           userID,
		WorkoutDayTypeID: dayTypeID,
		StartedAt:        time.Now(),
		Completed:        false,
	}
	err := s.workoutsRepo.Create(&workoutDay)
	if err != nil {
		fmt.Printf("%s: create workout error: %s\n", method, err.Error())
		return
	}

	previousWorkout, err := s.workoutsRepo.FindPreviousByType(userID, dayTypeID)
	if err != nil {
		err = s.createExercisesFromPresets(workoutDay.ID, dayTypeID)
	} else {
		err = s.createExercisesFromLastWorkout(workoutDay.ID, previousWorkout.ID)
	}
	if err != nil {
		fmt.Printf("%s: create exercises error: %s\n", method, err.Error())
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ <b>Тренировка создана!</b>\n\n")
	msg.ParseMode = constants.HtmlParseMode
	_, _ = tghelpers.SendMessage(s.bot, msg, method)

	s.showWorkoutProgress(chatID, workoutDay.ID)
}

func (s *serviceImpl) createExercisesFromPresets(workoutDayID, dayTypeID int64) error {
	method := "createExercisesFromPresets"
	fmt.Printf("%s: берем настройки количества повторений и веса из preset-ов\n", method)

	exercises := make([]models.Exercise, 0)
	dayType, err := s.dayTypesRepo.Get(dayTypeID)
	if err != nil {
		return err
	}

	fmt.Printf("%s: dayType: %d, preset: %s\n", method, dayType.ID, dayType.Preset)
	for index, exerciseType := range utils.SplitPreset(dayType.Preset) {
		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: exerciseType.ID,
			Index:          index,
		}
		for idx2, set := range exerciseType.Sets {
			newSet := models.Set{Index: idx2}
			if set.Minutes > 0 {
				newSet.Minutes = set.Minutes
			} else {
				newSet.Reps = set.Reps
				newSet.Weight = set.Weight
			}
			newExercise.Sets = append(newExercise.Sets, newSet)
		}
		exercises = append(exercises, newExercise)
	}

	return s.exercisesRepo.CreateBatch(exercises)
}

func (s *serviceImpl) createExercisesFromLastWorkout(workoutDayID, previousWorkoutID int64) error {
	method := "createExercisesFromLastWorkout"
	fmt.Printf("%s: берем настройки количества повторений и веса из последней тренировки: %d\n", method, previousWorkoutID)

	previousExercises, err := s.exercisesRepo.FindAllByWorkoutID(previousWorkoutID)
	if err != nil {
		return err
	}
	exercises := make([]models.Exercise, 0)
	for _, exercise := range previousExercises {
		newExercise := models.Exercise{
			WorkoutDayID:   workoutDayID,
			ExerciseTypeID: exercise.ExerciseTypeID,
			Index:          exercise.Index,
		}
		for _, set := range exercise.Sets {
			newSet := models.Set{
				Reps:    set.GetRealReps(),
				Weight:  set.GetRealWeight(),
				Minutes: set.GetRealMinutes(),
				Meters:  set.GetRealMeters(),
				Index:   set.Index,
			}
			newExercise.Sets = append(newExercise.Sets, newSet)
		}
		exercises = append(exercises, newExercise)
	}

	return s.exercisesRepo.CreateBatch(exercises)
}

func (s *serviceImpl) confirmDeleteWorkout(chatID int64, workoutID int64) {
	method := "confirmDeleteWorkout"
	workoutDay, err := s.workoutsRepo.Get(workoutID)
	if err != nil {
		return
	}

	dayType, err := s.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🗑️ *Удаление тренировки*\n\n"+
		"Вы уверены, что хотите удалить тренировку:\n"+
		"*%s*?\n\n"+
		"⚠️ Это действие нельзя отменить!", dayType.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("workout_delete_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("workout_show_progress_%d", workoutDay.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) deleteWorkout(chatID int64, workoutID int64) {
	method := "deleteWorkout"

	workoutDay, err := s.workoutsRepo.Get(workoutID)
	if err != nil {
		return
	}

	for _, exercise := range workoutDay.Exercises {
		deleteErr := s.setsRepo.DeleteAllBy(exercise.ID)
		if deleteErr != nil {
			return
		}
	}

	err = s.exercisesRepo.DeleteByWorkout(workoutID)
	if err != nil {
		return
	}

	err = s.workoutsRepo.Delete(&workoutDay)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Тренировка успешно удалена!")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.MyWorkouts, "workout_show_my"),
			tgbotapi.NewInlineKeyboardButtonData(messages.BackToMenu, "back_to_menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) startSpecificWorkout(chatID int64, workoutID int64) {
	method := "startSpecificWorkout"
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
		return
	}

	if workoutDay.Completed {
		msg := tgbotapi.NewMessage(chatID, "❌ Эта тренировка уже завершена. Создайте новую или повторите эту.")
		_, _ = tghelpers.SendMessage(s.bot, msg, method)
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

func (s *serviceImpl) confirmFinishWorkout(chatID int64, workoutDayID int64) {
	method := "confirmFinishWorkout"

	workoutDay, err := s.workoutsRepo.Get(workoutDayID)
	if err != nil {
		return
	}

	dayType, err := s.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🏁 *Завершение тренировки*\n\n"+
		"Вы уверены, что хотите завершить тренировку:\n"+
		"*%s*?\n\n"+
		"После завершения вы сможете просмотреть статистику, "+
		"но не сможете добавлять новые подходы.", dayType.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, завершить",
				fmt.Sprintf("workout_finish_%d", workoutDayID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, продолжить",
				fmt.Sprintf("workout_continue_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}

func (s *serviceImpl) finishWorkoutById(chatID int64, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	now := time.Now()
	workoutDay.Completed = true
	workoutDay.EndedAt = &now
	err := s.workoutsRepo.Save(&workoutDay)
	if err != nil {
		return
	}

	err = s.sessionsRepo.UpdateIsActive(workoutID, false)
	if err != nil {
		return
	}
	s.showWorkoutStatistics(chatID, workoutID)
}

func (s *serviceImpl) showWorkoutStatistics(chatID int64, workoutID int64) {
	method := "showWorkoutStatistics"
	text := s.statisticsService.ShowWorkoutStatistics(workoutID)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	_, _ = tghelpers.SendMessage(s.bot, msg, method)
}
