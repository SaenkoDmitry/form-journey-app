package service

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"

	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *serviceImpl) HandleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	fmt.Println("HandleMessage:", text)

	switch {
	case text == "🔙 В меню" || text == "/start" || text == "/menu":
		s.sendMainMenu(chatID, message.From)

	case text == "▶️ Начать тренировку" || text == "/start_workout":
		s.showWorkoutTypeMenu(chatID)

	case text == "📋 Мои тренировки" || text == "/workouts":
		s.showMyWorkouts(chatID, 0)

	case text == "📊 Статистика" || text == "/stats":
		s.showStatsMenu(chatID)

	case text == "⚙️ Настройки" || text == "/settings":
		s.settings(chatID)

	case text == "❓ Что умеет бот?" || text == "/about":
		s.about(chatID)

	default:
		s.handleState(chatID, text)
	}
}

func (s *serviceImpl) sendMainMenu(chatID int64, from *tgbotapi.User) {
	method := "sendMainMenu"

	text := "🏋️‍♂️ *Добро пожаловать в Бот для тренировок!* \n\n Выберите действие:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("▶️ Начать тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Мои тренировки"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
			tgbotapi.NewKeyboardButton("❓ Что умеет бот?"),
		),
	)
	keyboard.ResizeKeyboard = true

	s.createUserIfNotExists(chatID, from)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err := s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) createUserIfNotExists(chatID int64, from *tgbotapi.User) {
	_, err := s.usersRepo.GetByChatID(chatID)
	if err == nil {
		return
	}
	if errors.Is(err, users.NotFoundUserErr) {
		user, createErr := s.usersRepo.Create(chatID, from)
		if createErr != nil {
			return
		}

		// создаем дефолтную программу
		program, createErr := s.programsRepo.Create(user.ID, "#1 (базовая)")
		if createErr != nil {
			return
		}

		// прикрепляем программу к юзеру и сохраняем
		user.ActiveProgramID = &program.ID
		err = s.usersRepo.Save(user)
		if err != nil {
			return
		}
	}
}

func (s *serviceImpl) showWorkoutTypeMenu(chatID int64) {
	method := "showWorkoutTypeMenu"

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	program, err := s.programsRepo.Get(*user.ActiveProgramID)
	if err != nil {
		return
	}

	if len(program.DayTypes) == 0 {
		msg := tgbotapi.NewMessage(chatID, "Добавьте тренировочные дни в программу через '⚙️ Настройки'")
		msg.ParseMode = "Markdown"
		_, err = s.bot.Send(msg)
		handleErr(method, err)
		return
	}

	text := "*Выберите день тренировки:*"

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i, day := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(day.Name, fmt.Sprintf("create_workout_%d", day.ID)),
		)
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})

	//buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
	//	tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новый день", "create_day_type"),
	//)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "Markdown"
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) handleGetUserErr(chatID int64, method string, err error) {
	if errors.Is(err, users.NotFoundUserErr) {
		msg := tgbotapi.NewMessage(chatID, "Сначала создайте пользователя в боте, через команду /start")
		_, err = s.bot.Send(msg)
		handleErr(method, err)
	}
}

func handleErr(method string, err error) {
	if err != nil {
		fmt.Printf("\n %s: error is: %s \n", method, err.Error())
	}
}

const (
	showWorkoutsLimit = 4
)

func (s *serviceImpl) showMyWorkouts(chatID int64, offset int) {
	method := "showMyWorkouts"
	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	count, _ := s.workoutsRepo.Count(user.ID)

	limit := showWorkoutsLimit

	workouts, _ := s.workoutsRepo.Find(user.ID, offset, limit)

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

	var rows [][]tgbotapi.InlineKeyboardButton

	text := fmt.Sprintf("📋 *Ваши тренировки (%d):*\n\n", count)
	for i, workout := range workouts {
		status := "🟡"
		if workout.Completed {
			status = "✅"
			if workout.EndedAt != nil {
				status += fmt.Sprintf(" ~ %s",
					utils.BetweenTimes(workout.StartedAt, workout.EndedAt),
				)
			}
		}
		date := workout.StartedAt.Format("02.01.2006 15:04")

		dayType := workout.WorkoutDayType

		text += fmt.Sprintf("%d. *%s* %s\n   📅 %s\n\n",
			i+1+offset, dayType.Name, status, date)

		// buttons
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d", dayType.Name, i+1+offset),
				fmt.Sprintf("show_progress_%d", workout.ID)))
	}

	text += "Выберите тренировку для просмотра:"

	rows = append(rows, []tgbotapi.InlineKeyboardButton{})
	fmt.Println("offset", offset, "limit", limit, "count", count)
	if offset >= limit {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("⬅️ Предыдущие",
			fmt.Sprintf("my_workouts_%d", offset-limit)))
	}
	if offset+limit < int(count) {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("➡️ Следующие",
			fmt.Sprintf("my_workouts_%d", offset+limit)))
	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("🔙 В начало", "my_workouts"))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) showStatsMenu(chatID int64) {
	method := "showStatsMenu"
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
	_, err := s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) changeProgram(chatID, programID int64) {
	method := "changeProgram"
	fmt.Sprintf("%s: programID: %d", method, programID)

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	user.ActiveProgramID = &programID
	err = s.usersRepo.Save(user)
	if err != nil {
		fmt.Printf("%s: %s\n", method, err.Error())
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Успешно обновлено!")
	_, err = s.bot.Send(msg)
	handleErr(method, err)
	s.settings(chatID)
}

func (s *serviceImpl) confirmDeleteProgram(chatID, programID int64) {
	method := "confirmDeleteProgram"

	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	text := fmt.Sprintf("🗑️ *Удаление тренировки*\n\n"+
		"Вы уверены, что хотите удалить тренировку:\n"+
		"*%s*?\n\n"+
		"❌ Это действие нельзя отменить!", program.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("delete_program_%d", programID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("edit_program_%d", programID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) deleteProgram(chatID, programID int64) {
	method := "deleteProgram"

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	if user.ActiveProgramID == &programID {
		msg := tgbotapi.NewMessage(chatID, "Нельзя удалить текущую программу 😓")
		_, err = s.bot.Send(msg)
		return
	}

	program, err := s.programsRepo.Get(programID)
	if err != nil {
		return
	}

	err = s.programsRepo.Delete(&program)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, "✅ Успешно удалено!")
	_, err = s.bot.Send(msg)
	handleErr(method, err)
	s.settings(chatID)
}

func (s *serviceImpl) settings(chatID int64) {
	method := "settings"

	user, err := s.usersRepo.GetByChatID(chatID)
	if err != nil {
		s.handleGetUserErr(chatID, method, err)
		return
	}

	programs, err := s.programsRepo.FindAll(user.ID)
	if err != nil {
		return
	}

	addNewProgram := tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новую", "create_program")

	if len(programs) == 0 {
		msg := tgbotapi.NewMessage(chatID, "🥲 У вас нет тренировочных программ, создайте первую!")
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(addNewProgram))
		_, err = s.bot.Send(msg)
		handleErr(method, err)
		return
	}

	text := &bytes.Buffer{}
	text.WriteString("*Ваши программы:*\n\n")

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, program := range programs {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}

		if &program.ID == user.ActiveProgramID {
			text.WriteString(fmt.Sprintf("• %s %s *(текущая)*\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		} else {
			text.WriteString(fmt.Sprintf("• %s %s\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		}

		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(program.Name, fmt.Sprintf("edit_program_%d", program.ID)))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(addNewProgram))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err = s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) about(chatID int64) {
	method := "about"
	msg := tgbotapi.NewMessage(chatID, `
	*Цель бота*: помощь в учете тренировок, отслеживании весов / повторов, установка таймеров, просмотр статистики

	# *Что умеет бот?*

	*1).* В пункте меню *'▶️ Начать тренировку'* есть следующие функции:
		
		• ⚠️ пока в боте задана лишь *одна (!)* тренировочная программа (в дальнейшем появится возможность составлять свои)
		
		• ✍️ бот позволяет записывать запланированные/фактические веса и повторы
		
		• 📕 бот умеет запоминать веса/повторы с прошлой тренировки
		
		• ⏱️ бот умеет засекать время на таймере между подходами
		
		• 🤓 бот предоставляет видео с техникой выполнения упражнения

	*2).* 📖 В пункте меню *'📋 Мои тренировки'* можно посмотреть историю своих тренировок

	*3).* 🏆 В пункте меню *'📊 Статистика'* можно посмотреть краткую сводку тренировок
		• кол-во за период
		• среднее время силовых тренировок
		• отдельно время кардио тренировок
		• вышеперечисленное в разрезе: неделя, месяц, общая
	`)

	msg.ParseMode = "Markdown"
	_, err := s.bot.Send(msg)
	handleErr(method, err)
}

func (s *serviceImpl) handleState(chatID int64, text string) {
	method := "handleState"
	state, exists := s.userStates[chatID]
	if !exists {
		return
	}

	var err error

	switch {
	case strings.HasPrefix(state, "awaiting_reps_"):
		err = s.awaitingEnterData(
			chatID, state,
			func() (interface{}, error) { return strconv.ParseInt(text, 10, 64) },
			func(nextSet models.Set, value interface{}) models.Set {
				reps, ok := value.(int64)
				if !ok {
					return models.Set{}
				}
				nextSet.FactReps = int(reps)
				if int(reps) != nextSet.Reps {
					nextSet.FactReps = int(reps)
				} else {
					nextSet.FactReps = 0
				}
				return nextSet
			},
			"❌ Неверный формат числа повторений. Введите целое число (например: 42)",
			"✅ Количество повторений обновлено",
		)
	case strings.HasPrefix(state, "awaiting_weight_"):
		err = s.awaitingEnterData(
			chatID, state,
			func() (interface{}, error) { return strconv.ParseFloat(text, 32) },
			func(nextSet models.Set, value interface{}) models.Set {
				weight, ok := value.(float64)
				if !ok {
					return models.Set{}
				}
				if float32(weight) != nextSet.Weight {
					nextSet.FactWeight = float32(weight)
				} else {
					nextSet.FactWeight = float32(0)
				}
				return nextSet
			},
			"❌ Неверный формат веса. Введите число (например: 42.5)",
			"✅ Вес обновлен",
		)

	case strings.HasPrefix(state, "awaiting_minutes_"):
		err = s.awaitingEnterData(
			chatID, state,
			func() (interface{}, error) { return strconv.ParseInt(text, 10, 64) },
			func(nextSet models.Set, value interface{}) models.Set {
				minutes, ok := value.(int64)
				if !ok {
					return models.Set{}
				}
				if int(minutes) != nextSet.Minutes {
					nextSet.FactMinutes = int(minutes)
				} else {
					nextSet.FactMinutes = 0
				}
				return nextSet
			},
			"❌ Неверный формат минут. Введите число (например: 42)",
			"✅ Время обновлено",
		)

	case strings.HasPrefix(state, "awaiting_program_name_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_program_name_"), 10, 64)
		program, err := s.programsRepo.Get(programID)
		if err != nil {
			return
		}
		program.Name = text
		err = s.programsRepo.Save(&program)
		if err != nil {
			return
		}
		s.settings(chatID)

	case strings.HasPrefix(state, "awaiting_day_preset_"):
		preset := text
		if !utils.IsValidPreset(preset) {
			msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат. Введите снова (например: 15*100,10*150)")
			s.bot.Send(msg)
			fmt.Printf("%s: invalid preset: %s\n", method, preset)
			return
		}
		parts := strings.Split(strings.TrimPrefix(state, "awaiting_day_preset_"), "_")
		if len(parts) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(parts[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(parts[1], 10, 64)

		var dayType models.WorkoutDayType
		dayType, err = s.dayTypesRepo.Get(dayTypeID)
		if err != nil {
			return
		}
		if dayType.Preset != "" {
			dayType.Preset += ";"
		}
		dayType.Preset += fmt.Sprintf("%d:[%s]", exerciseTypeID, preset)
		err = s.dayTypesRepo.Save(&dayType)
		if err != nil {
			return
		}
		s.editProgram(chatID, dayType.WorkoutProgramID)

	case strings.HasPrefix(state, "awaiting_day_name_for_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_day_name_for_program_"), 10, 64)

		dayType, createErr := s.dayTypesRepo.Create(&models.WorkoutDayType{
			WorkoutProgramID: programID,
			Name:             text,
			CreatedAt:        time.Now(),
		})
		if createErr != nil {
			return
		}
		s.addNewDayTypeExercise(chatID, dayType.ID)
	}

	handleErr(method, err)
}

func (s *serviceImpl) awaitingEnterData(
	chatID int64,
	state string,
	parseValue func() (interface{}, error),
	handleSet func(s models.Set, result interface{}) models.Set,
	formatMsg, successMsg string,
) error {
	parts := strings.Split(state, "_")
	if len(parts) < 3 {
		return errors.New("incorrect input")
	}
	exerciseID, _ := strconv.ParseInt(parts[2], 10, 64)

	result, err := parseValue()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, formatMsg)
		_, err = s.bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}

	exercise, _ := s.exercisesRepo.Get(exerciseID)
	nextSet := exercise.NextSet()

	if nextSet.ID != 0 {
		nextSet = handleSet(nextSet, result)
		err = s.setsRepo.Save(&nextSet)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(chatID, successMsg)
		if _, err = s.bot.Send(msg); err != nil {
			return err
		}
	}
	s.userStates[chatID] = ""
	s.showCurrentExerciseSession(chatID, exercise.WorkoutDayID)
	return nil
}
