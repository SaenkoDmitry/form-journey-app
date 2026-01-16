package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"strconv"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	fmt.Println("HandleCallback:", data)

	switch {
	case data == "back_to_menu":
		s.sendMainMenu(chatID, callback.From)

	// programs
	case strings.HasPrefix(data, "create_program"):
		s.createProgram(chatID)

	case strings.HasPrefix(data, "edit_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "edit_program_"), 10, 64)
		s.editProgram(chatID, programID)

	case strings.HasPrefix(data, "change_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_program_"), 10, 64)
		s.changeProgram(chatID, programID)

	case strings.HasPrefix(data, "change_name_of_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_name_of_program_"), 10, 64)
		s.askForNewProgramName(chatID, programID)

	case strings.HasPrefix(data, "confirm_delete_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "confirm_delete_program_"), 10, 64)
		s.confirmDeleteProgram(chatID, programID)

	case strings.HasPrefix(data, "delete_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "delete_program_"), 10, 64)
		s.deleteProgram(chatID, programID)

	case strings.HasPrefix(data, "edit_day_type_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "edit_day_type_"), 10, 64)
		s.addNewDayTypeExercise(chatID, dayTypeID)

	// workouts
	case strings.HasPrefix(data, "create_workout_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "create_workout_"), 10, 64)
		s.createWorkoutDay(chatID, dayTypeID)

	case strings.HasPrefix(data, "start_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "start_workout_"), 10, 64)
		s.startSpecificWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "unpin_and_cancel_timer_"):
		timerID := strings.TrimPrefix(data, "unpin_and_cancel_timer_")
		s.unpinAndCancelTimer(chatID, timerID)

	case strings.HasPrefix(data, "my_workouts"):
		if data == "my_workouts" {
			s.showMyWorkouts(chatID, 0)
			return
		}
		offset, _ := strconv.ParseInt(strings.TrimPrefix(data, "my_workouts_"), 10, 64)
		s.showMyWorkouts(chatID, int(offset))

	case strings.HasPrefix(data, "confirm_delete_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "confirm_delete_workout_"), 10, 64)
		s.confirmDeleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "delete_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "delete_workout_"), 10, 64)
		s.deleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "continue_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "continue_workout_"), 10, 64)
		s.showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "show_progress_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "show_progress_"), 10, 64)
		s.showWorkoutProgress(chatID, workoutID)

	case strings.HasPrefix(data, "finish_workout_id_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "finish_workout_id_"), 10, 64)
		s.confirmFinishWorkout(chatID, workoutDayID)

	case strings.HasPrefix(data, "do_finish_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "do_finish_workout_"), 10, 64)
		s.finishWorkoutById(chatID, workoutDayID)

	case strings.HasPrefix(data, "stats_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "stats_workout_"), 10, 64)
		s.showWorkoutStatistics(chatID, workoutID)

	// timer
	case strings.HasPrefix(data, "start_timer_"):
		fmt.Println("start_timer_: data: ", data)
		parts := strings.Split(data, "_")
		if len(parts) >= 5 && parts[3] == "ex" {
			seconds, _ := strconv.Atoi(parts[2])
			exerciseID, _ := strconv.ParseInt(parts[4], 10, 64)
			s.startRestTimerWithExercise(chatID, seconds, exerciseID)
		}

	// exercises
	case strings.HasPrefix(data, "complete_set_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "complete_set_ex_"), 10, 64)
		s.completeExerciseSet(chatID, exerciseID)

	case strings.HasPrefix(data, "add_one_more_set_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "add_one_more_set_"), 10, 64)
		s.addOneMoreSet(chatID, exerciseID)

	case strings.HasPrefix(data, "prev_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "prev_exercise_"), 10, 64)
		s.moveToPrevExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "show_current_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "show_current_exercise_"), 10, 64)
		s.showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "next_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "next_exercise_"), 10, 64)
		s.moveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "show_exercise_hint_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "show_exercise_hint_"), 10, 64)
		s.showExerciseHint(chatID, workoutID)

	case strings.HasPrefix(data, "add_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "add_exercise_"), 10, 64)
		s.addExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "select_exercise_"):
		text := strings.TrimPrefix(data, "select_exercise_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutDayID, _ := strconv.ParseInt(arr[0], 10, 64)
			code := arr[1]
			fmt.Println("workoutID:", workoutDayID, "code:", code)
			s.selectExercise(chatID, workoutDayID, code)
		}

	case strings.HasPrefix(data, "add_specific_exercise_"):
		text := strings.TrimPrefix(data, "add_specific_exercise_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutID, _ := strconv.ParseInt(arr[0], 10, 64)
			internalExerciseID, _ := strconv.ParseInt(arr[1], 10, 64)
			fmt.Println("workoutID:", workoutID, "internalExerciseID:", internalExerciseID)
			s.addSpecificExercise(chatID, workoutID, internalExerciseID)
		}

	case strings.HasPrefix(data, "confirm_delete_exercise_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "confirm_delete_exercise_"), 10, 64)
		s.confirmDeleteExercise(chatID, exerciseID)

	case strings.HasPrefix(data, "delete_exercise_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "delete_exercise_"), 10, 64)
		s.deleteExercise(chatID, exerciseID)

	// settings
	case strings.HasPrefix(data, "change_reps_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_reps_ex_"), 10, 64)
		s.askForNewReps(chatID, exerciseID)

	case strings.HasPrefix(data, "change_weight_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_weight_ex_"), 10, 64)
		s.askForNewWeight(chatID, exerciseID)

	case strings.HasPrefix(data, "change_minutes_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_minutes_ex_"), 10, 64)
		s.askForNewMinutes(chatID, exerciseID)

	case strings.HasPrefix(data, "change_meters_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_meters_ex_"), 10, 64)
		s.askForNewMeters(chatID, exerciseID)

	case strings.HasPrefix(data, "create_day_type_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "create_day_type_"), 10, 64)
		s.askForNewDayName(chatID, programID)

	case strings.HasPrefix(data, "day_type_select_exercise_"):
		text := strings.Split(strings.TrimPrefix(data, "day_type_select_exercise_"), "_")
		if len(text) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(text[0], 10, 64)
		exerciseGroupCode := text[1]
		s.addForDaySpecificExercise(chatID, dayTypeID, exerciseGroupCode)

	case strings.HasPrefix(data, "day_type_add_specific_exercise_"):
		text := strings.Split(strings.TrimPrefix(data, "day_type_add_specific_exercise_"), "_")
		if len(text) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(text[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(text[1], 10, 64)

		s.askForPreset(chatID, dayTypeID, exerciseTypeID)

	// stats
	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		s.showStatistics(chatID, period)
	}
}

func (s *serviceImpl) showWorkoutProgress(chatID, workoutID int64) {
	workoutDay, _ := s.workoutsRepo.Get(workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		s.bot.Send(msg)
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

	session, _ := s.sessionsRepo.GetByWorkoutID(workoutID)

	toWorkoutButton := tgbotapi.NewInlineKeyboardButtonData("▶️ К тренировке", fmt.Sprintf("show_current_exercise_%d", workoutID))
	if session.ID == 0 {
		toWorkoutButton = tgbotapi.NewInlineKeyboardButtonData("▶️ Начать", fmt.Sprintf("start_workout_%d", workoutDay.ID))
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	if !workoutDay.Completed {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Добавить упражнение", fmt.Sprintf("add_exercise_%d", workoutID)),
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("confirm_delete_workout_%d", workoutID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				toWorkoutButton,
			),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📊 Статистика", fmt.Sprintf("stats_workout_%d", workoutID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "my_workouts"),
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить",
					fmt.Sprintf("confirm_delete_workout_%d", workoutID)),
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Html"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) createProgram(chatID int64) {
	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	programs, err := s.programsRepo.FindAll(user.ID)
	if err != nil {
		return
	}

	_, err = s.programsRepo.Create(user.ID, fmt.Sprintf("#%d", len(programs)+1))
	if err != nil {
		return
	}

	s.settings(chatID)
}

func (s *serviceImpl) editProgram(chatID int64, programID int64) {
	_, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	text := &bytes.Buffer{}
	text.WriteString(fmt.Sprintf("*Программа: %s*\n\n", program.Name))
	text.WriteString("*Список дней:*\n\n")
	for i, dayType := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(dayType.Name, fmt.Sprintf("edit_day_type_%d", dayType.ID)),
		)

		text.WriteString(fmt.Sprintf("*%d. %s*\n", i+1, dayType.Name))
		text.WriteString(fmt.Sprintf("%s \n\n", s.formatPreset(dayType.Preset)))
	}
	text.WriteString("*Выберите день, в который хотите добавить упражнения:*")

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить день", fmt.Sprintf("create_day_type_%d", programID)),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Переименовать", fmt.Sprintf("change_name_of_program_%d", programID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👑 Выбрать текущей", fmt.Sprintf("change_program_%d", programID)),
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", fmt.Sprintf("delete_program_%d", programID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) formatPreset(preset string) string {
	exercises := utils.SplitPreset(preset)

	buffer := &bytes.Buffer{}
	for _, ex := range exercises {

		exerciseType, err := s.exerciseTypesRepo.Get(ex.ID)
		if err != nil {
			continue
		}
		buffer.WriteString(fmt.Sprintf("• *%s*\n", exerciseType.Name))
		buffer.WriteString(fmt.Sprintf("    • "))
		for i, set := range ex.Sets {
			if i > 0 {
				buffer.WriteString(", ")
			}
			if set.Minutes > 0 {
				buffer.WriteString(fmt.Sprintf("%d мин", set.Minutes))
			} else {
				buffer.WriteString(fmt.Sprintf("%d \\* %.0f кг", set.Reps, set.Weight))
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (s *serviceImpl) createWorkoutDay(chatID int64, dayTypeID int64) {
	method := "createWorkoutDay"
	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	workoutDay := models.WorkoutDay{
		UserID:           user.ID,
		WorkoutDayTypeID: dayTypeID,
		StartedAt:        time.Now(),
		Completed:        false,
	}
	err = s.workoutsRepo.Create(&workoutDay)
	if err != nil {
		fmt.Printf("%s: create workout error: %s\n", method, err.Error())
		return
	}

	previousWorkout, err := s.workoutsRepo.FindPreviousByType(user.ID, dayTypeID)
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
	msg.ParseMode = "Html"
	s.bot.Send(msg)

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
		"❌ Это действие нельзя отменить!", dayType.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("delete_workout_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("show_progress_%d", workoutDay.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) confirmDeleteExercise(chatID int64, exerciseID int64) {
	exercise, _ := s.exercisesRepo.Get(exerciseID)

	exerciseObj, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🗑️ *Удаление упражнения из тренировочного дня*\n\n"+
		"Вы уверены, что хотите удалить упражнение:\n"+
		"*%s*?\n\n"+
		"❌ Это действие нельзя отменить!", exerciseObj.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("delete_exercise_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("start_workout_%d", exercise.WorkoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) deleteExercise(chatID int64, exerciseID int64) {
	exercise, err := s.exercisesRepo.Get(exerciseID)
	if err != nil {
		return
	}

	err = s.exercisesRepo.Delete(exerciseID)
	if err != nil {
		return
	}
	// session, _ := s.sessionsRepo.GetByWorkoutID(exercise.WorkoutDayID)
	// session.CurrentExerciseIndex++
	// s.sessionsRepo.Save(&session)

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
}

func (s *serviceImpl) deleteWorkout(chatID int64, workoutID int64) {
	method := "deleteWorkout"

	workoutDay, err := s.workoutsRepo.Get(workoutID)
	if err != nil {
		return
	}

	for _, exercise := range workoutDay.Exercises {
		deleteErr := s.setsRepo.Delete(exercise.ID)
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
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
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

	exerciseObj, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Упражнение не найдено.")
		s.bot.Send(msg)
		return
	}

	dayType, err := s.dayTypesRepo.Get(workoutDay.WorkoutDayTypeID)
	if err != nil {
		return
	}

	text.WriteString(fmt.Sprintf("<b>%s</b>\n\n", dayType.Name))
	text.WriteString(fmt.Sprintf("<b>Упражнение %d/%d:</b> %s\n\n", exerciseIndex+1, len(workoutDay.Exercises), exerciseObj.Name))
	if exerciseObj.Accent != "" {
		text.WriteString(fmt.Sprintf("<b>Акцент:</b> %s\n\n", exerciseObj.Accent))
	}

	for _, set := range exercise.Sets {
		text.WriteString(set.String(workoutDay.Completed))
	}

	var changeSettingsButtons []tgbotapi.InlineKeyboardButton
	if len(exercise.Sets) > 0 && exercise.Sets[0].Minutes > 0 {
		changeSettingsButtons = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Minutes, fmt.Sprintf("change_minutes_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Meters > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Meters, fmt.Sprintf("change_meters_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Reps > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Reps, fmt.Sprintf("change_reps_ex_%d", exercise.ID)),
		)
	}

	if len(exercise.Sets) > 0 && exercise.Sets[0].Weight > 0 {
		changeSettingsButtons = append(changeSettingsButtons,
			tgbotapi.NewInlineKeyboardButtonData(messages.Weight, fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Set, fmt.Sprintf("complete_set_ex_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Timer, fmt.Sprintf("start_timer_%d_ex_%d", exercise.ExerciseType.RestInSeconds, exercise.ID)),
		),
		changeSettingsButtons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.AddSet, fmt.Sprintf("add_one_more_set_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить", fmt.Sprintf("finish_workout_id_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Удалить", fmt.Sprintf("confirm_delete_exercise_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧐 Показать технику", fmt.Sprintf("show_exercise_hint_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("📊 Прогресс", fmt.Sprintf("show_progress_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Пред", fmt.Sprintf("prev_exercise_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("➡️ След", fmt.Sprintf("next_exercise_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Html"
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	if err != nil {
		fmt.Println("showCurrentExerciseSession", err.Error())
	}
}

func (s *serviceImpl) addOneMoreSet(chatID int64, exerciseID int64) {
	exercise, err := s.exercisesRepo.Get(exerciseID)
	if err != nil || len(exercise.Sets) == 0 {
		return
	}
	lastSet := exercise.Sets[len(exercise.Sets)-1]
	err = s.setsRepo.Save(&models.Set{
		ExerciseID: exercise.ID,
		Reps:       lastSet.Reps,
		Weight:     lastSet.Weight,
		Minutes:    lastSet.Minutes,
		Meters:     lastSet.Meters,
		Index:      lastSet.Index + 1,
	})
	if err != nil {
		fmt.Println("cannot create set:", err.Error())
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ <b>Еще один подход добавлен!</b>")
	msg.ParseMode = "Html"
	s.bot.Send(msg)

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
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

	text := fmt.Sprintf("✅ *Подход завершен!*\n\n")
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)

	exerciseType, err := s.exerciseTypesRepo.Get(exercise.ExerciseTypeID)
	if err != nil {
		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
		return
	}

	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	s.startRestTimerWithExercise(chatID, exerciseType.RestInSeconds, exerciseID)
}

func (s *serviceImpl) unpinAndCancelTimer(chatID int64, timerID string) {
	s.timerStore.StopTimer(chatID, timerID)
}

func (s *serviceImpl) startRestTimerWithExercise(chatID int64, seconds int, exerciseID int64) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(messages.RestTimer, seconds))
	newTimerID := s.timerStore.NewTimer(chatID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить", fmt.Sprintf("unpin_and_cancel_timer_%s", newTimerID)),
		),
	)
	msg.ReplyMarkup = keyboard

	var message tgbotapi.Message
	message, _ = s.bot.Send(msg)
	s.pinMessage(chatID, message)

	go func() {
		remaining := seconds

		for remaining > 0 {
			time.Sleep(1 * time.Second)
			remaining--
			if !s.timerStore.HasTimer(chatID, newTimerID) {
				s.unpinMessage(chatID, message)
				editMsg := tgbotapi.NewEditMessageText(chatID, message.MessageID, "Таймер отменен")
				s.bot.Send(editMsg)
				return
			}

			var err error
			if remaining%10 == 0 || remaining <= 20 {
				editMsg := tgbotapi.NewEditMessageTextAndMarkup(chatID, message.MessageID,
					fmt.Sprintf(messages.RestTimer, remaining), keyboard)
				if message, err = s.bot.Send(editMsg); err != nil {
					fmt.Println("cannot edit msg")
				}
			}
		}

		editMsg := tgbotapi.NewEditMessageText(
			chatID,
			message.MessageID,
			"🔔 *Время отдыха закончилось!*\n\n Приступайте к следующему подходу! 💪",
		)
		editMsg.ParseMode = "Markdown"
		editMessage, _ := s.bot.Send(editMsg)
		s.unpinMessage(chatID, editMessage)

		exercise, _ := s.exercisesRepo.Get(exerciseID)

		s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	}()
}

func (s *serviceImpl) pinMessage(chatID int64, message tgbotapi.Message) {
	pinChatMessageConfig := tgbotapi.PinChatMessageConfig{
		ChatID:              chatID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	if _, err := s.bot.Request(pinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
}

func (s *serviceImpl) unpinMessage(chatID int64, message tgbotapi.Message) {
	unpinChatMessageConfig := tgbotapi.UnpinChatMessageConfig{
		ChatID:    chatID,
		MessageID: message.MessageID,
	}
	if _, err := s.bot.Request(unpinChatMessageConfig); err != nil {
		fmt.Println("cannot pin message:", message.MessageID)
	}
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
				"Хотите завершить тренировку или добавить еще упражнения?")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
					fmt.Sprintf("finish_workout_id_%d", workoutID)),
				tgbotapi.NewInlineKeyboardButtonData("➕ Добавить упражнение",
					fmt.Sprintf("add_exercise_%d", workoutID)),
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

func (s *serviceImpl) showExerciseHint(chatID int64, workoutID int64) {
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

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("show_current_exercise_%d", workoutID)),
	))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, utils.WrapYandexLink(exercise.ExerciseType.Url))
	msg.ParseMode = "Html"
	msg.ReplyMarkup = keyboard
	_, err := s.bot.Send(msg)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
}

func (s *serviceImpl) addExercise(chatID int64, workoutID int64) {
	text := "*Выберите группу мышц:*"

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	groups, err := s.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return
	}

	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1], tgbotapi.NewInlineKeyboardButtonData(group.Name,
			fmt.Sprintf("select_exercise_%d_%s", workoutID, group.Code)))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) selectExercise(chatID int64, workoutID int64, exerciseGroupCode string) {
	group, err := s.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return
	}

	text := fmt.Sprintf("*Тип: %s \n\n Выберите упражнение из списка:*", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	exerciseTypes, err := s.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return
	}

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("add_specific_exercise_%d_%d", workoutID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}

func (s *serviceImpl) addSpecificExercise(chatID int64, workoutID int64, exerciseTypeID int64) {
	fmt.Println("addSpecificExercise:", "workoutID:", workoutID, "exerciseTypeID:", exerciseTypeID)

	exerciseObj, err := s.exerciseTypesRepo.Get(exerciseTypeID)
	if err != nil {
		return
	}

	fmt.Println("newExercise:", exerciseObj)

	workout, _ := s.workoutsRepo.Get(workoutID)
	idx := 0
	if len(workout.Exercises) > 0 {
		lastExercise := workout.Exercises[len(workout.Exercises)-1]
		idx = lastExercise.Index + 1
	}
	newExercise := models.Exercise{
		ExerciseTypeID: exerciseObj.ID,
		Index:          idx,
		WorkoutDayID:   workoutID,
		Sets: []models.Set{
			{Index: 1}, // по дефолту один подход
		},
	}
	exerciseTypeObj, err := s.exerciseTypesRepo.Get(newExercise.ExerciseTypeID)
	if err == nil {
		if strings.Contains(exerciseTypeObj.Units, "meters") {
			newExercise.Sets[0].Meters = 100
		}
		if strings.Contains(exerciseTypeObj.Units, "minutes") {
			newExercise.Sets[0].Minutes = 1
		}
		if strings.Contains(exerciseTypeObj.Units, "reps") {
			newExercise.Sets[0].Reps = 10
		}
		if strings.Contains(exerciseTypeObj.Units, "weight") {
			newExercise.Sets[0].Weight = 10
		}
	}
	workout.Exercises = append(workout.Exercises, newExercise)

	s.workoutsRepo.Save(&workout)

	msg := tgbotapi.NewMessage(chatID, "Упражнение добавлено! ✅")
	msg.ParseMode = "Markdown"
	s.bot.Send(msg)

	s.showWorkoutProgress(chatID, workoutID)
}

func (s *serviceImpl) addForDaySpecificExercise(chatID int64, dayTypeID int64, exerciseGroupCode string) {
	group, err := s.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return
	}

	text := fmt.Sprintf("*Тип: %s \n\n Выберите упражнение из списка:*", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	exerciseTypes, err := s.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return
	}

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("day_type_add_specific_exercise_%d_%d", dayTypeID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)

	//msg := tgbotapi.NewMessage(chatID, "Упражнение добавлено! ✅")
	//msg.ParseMode = "Markdown"
	//_, err = s.bot.Send(msg)
	//handleErr(method, err)
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
				fmt.Sprintf("do_finish_workout_%d", workoutDayID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, продолжить",
				fmt.Sprintf("continue_workout_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
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
	text := s.statisticsService.ShowWorkoutStatistics(workoutID)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) showStatistics(chatID int64, period string) {
	method := "showStatistics"
	user, err := s.GetUserByChatID(chatID)
	if err != nil {
		return
	}

	text := s.statisticsService.ShowStatistics(user.ID, period)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Html"
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) GetUserByChatID(chatID int64) (*models.User, error) {
	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, users.NotFoundUserErr) {
			msg := tgbotapi.NewMessage(chatID, "Сначала создайте пользователя в боте, через команду /start")
			_, err = s.bot.Send(msg)
			if err != nil {
				fmt.Printf("Error is: %v\n", err)
			}
		}
	}
	return user, nil
}

func (s *serviceImpl) askForNewReps(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_reps_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewReps)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewWeight(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_weight_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewWeight)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewMinutes(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_minutes_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewTime)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewMeters(chatID int64, exerciseID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_meters_%d", exerciseID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewMeters)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewDayName(chatID, programID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_name_for_program_%d", programID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterWorkoutDayName)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForNewProgramName(chatID, programID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_program_name_%d", programID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterNewProgramName)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) askForPreset(chatID, dayTypeID, exerciseTypeID int64) {
	s.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_preset_%d_%d", dayTypeID, exerciseTypeID))
	msg := tgbotapi.NewMessage(chatID, messages.EnterPreset)
	msg.ParseMode = "Html"
	s.bot.Send(msg)
}

func (s *serviceImpl) addNewDayTypeExercise(chatID, dayTypeID int64) {
	text := messages.SelectGroupOfMuscle

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	groups, err := s.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return
	}

	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(group.Name, fmt.Sprintf("day_type_select_exercise_%d_%s", dayTypeID, group.Code)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Html"
	msg.ReplyMarkup = keyboard
	s.bot.Send(msg)
}
