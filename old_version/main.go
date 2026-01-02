package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ================== Модели данных ==================
type User struct {
	ID        int64 `gorm:"primaryKey"`
	Username  string
	ChatID    int64
	CreatedAt time.Time
}

type WorkoutDay struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Name      string
	Exercises []Exercise `gorm:"foreignKey:WorkoutDayID"`
	StartedAt time.Time
	EndedAt   *time.Time
	Completed bool
}

type Exercise struct {
	ID           int64 `gorm:"primaryKey"`
	WorkoutDayID int64
	Name         string
	Sets         []Set `gorm:"foreignKey:ExerciseID"`
	TargetSets   int
	TargetReps   int
}

type Set struct {
	ID          int64 `gorm:"primaryKey"`
	ExerciseID  int64
	Reps        int
	Weight      float32
	Completed   bool
	CompletedAt *time.Time
}

type WorkoutSession struct {
	ID                   int64 `gorm:"primaryKey"`
	WorkoutDayID         int64
	CurrentExerciseIndex int
	StartedAt            time.Time
	IsActive             bool
}

type UserSetting struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Key       string
	Value     string
	UpdatedAt time.Time
}

// ================== Конфигурация ==================
type Config struct {
	TelegramToken string `json:"telegram_token"`
}

var (
	bot        *tgbotapi.BotAPI
	db         *gorm.DB
	userStates = make(map[int64]string)
)

// ================== Главная функция ==================
func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Config file not found")
	}
	defer configFile.Close()

	var config Config
	json.NewDecoder(configFile).Decode(&config)

	bot, err = tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	db, err = gorm.Open(sqlite.Open("workout_bot.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect database")
	}

	db.AutoMigrate(&User{}, &WorkoutDay{}, &Exercise{}, &Set{}, &WorkoutSession{}, &UserSetting{})

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(update.CallbackQuery)
		}
	}
}

// ================== Обработчики ==================
func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text

	user, _ := getUser(chatID, message.From.UserName)

	switch {
	case text == "/start" || text == "/menu" || text == "🔙 В меню":
		sendMainMenu(chatID)

	case text == "/new_workout" || text == "➕ Создать тренировку":
		showWorkoutTypeMenu(chatID)

	case text == "/start_workout" || text == "▶️ Начать тренировку":
		startActiveWorkout(chatID, user.ID)

	case text == "/stats" || text == "📊 Статистика":
		showStatsMenu(chatID, user.ID)

	case text == "⚙️ Настройки":
		showSettingsMenu(chatID)

	case text == "📋 Мои тренировки" || text == "/workouts":
		showMyWorkouts(chatID)

	default:
		handleState(chatID, user.ID, text)
	}
}

func handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch {
	case strings.HasPrefix(data, "create_workout_"):
		workoutType := strings.TrimPrefix(data, "create_workout_")
		createWorkoutDay(chatID, workoutType)

	case strings.HasPrefix(data, "start_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "start_workout_"), 10, 64)
		startSpecificWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "start_active_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "start_active_workout_"), 10, 64)
		startSpecificWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "view_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "view_workout_"), 10, 64)
		showWorkoutDetails(chatID, workoutID)

	case strings.HasPrefix(data, "edit_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "edit_workout_"), 10, 64)
		editWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "add_exercise_to_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "add_exercise_to_"), 10, 64)
		askForNewExercise(chatID, workoutID)

	case strings.HasPrefix(data, "confirm_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "confirm_delete_"), 10, 64)
		confirmDeleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "delete_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "delete_workout_"), 10, 64)
		deleteWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "stats_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "stats_workout_"), 10, 64)
		showWorkoutStatistics(chatID, workoutID)

	case strings.HasPrefix(data, "repeat_workout_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "repeat_workout_"), 10, 64)
		repeatWorkout(chatID, workoutID)

	case strings.HasPrefix(data, "complete_set_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "complete_set_ex_"), 10, 64)
		completeExerciseSet(chatID, exerciseID)

	case strings.HasPrefix(data, "add_reps_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "add_reps_ex_"), 10, 64)
		addRepsToLastSet(chatID, exerciseID)

	case strings.HasPrefix(data, "change_weight_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_weight_ex_"), 10, 64)
		askForNewWeight(chatID, exerciseID)

	case strings.HasPrefix(data, "rest_timer_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "rest_timer_ex_"), 10, 64)
		showRestTimerOptions(chatID, exerciseID)

	case strings.HasPrefix(data, "next_exercise_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "next_exercise_"), 10, 64)
		moveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "next_exercise_wd_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "next_exercise_wd_"), 10, 64)
		moveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "show_progress_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "show_progress_"), 10, 64)
		showWorkoutProgress(chatID, workoutDayID)

	case strings.HasPrefix(data, "finish_workout_id_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "finish_workout_id_"), 10, 64)
		confirmFinishWorkout(chatID, workoutDayID)

	case strings.HasPrefix(data, "do_finish_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "do_finish_workout_"), 10, 64)
		finishWorkoutById(chatID, workoutDayID)

	case strings.HasPrefix(data, "continue_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "continue_workout_"), 10, 64)
		showCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "timer_"):
		parts := strings.Split(data, "_")
		if len(parts) >= 2 {
			seconds, _ := strconv.Atoi(parts[1])
			if len(parts) >= 4 && parts[2] == "ex" {
				exerciseID, _ := strconv.ParseInt(parts[3], 10, 64)
				startRestTimerWithExercise(chatID, seconds, exerciseID)
			} else {
				startRestTimer(chatID, seconds)
			}
		}

	case strings.HasPrefix(data, "start_timer_"):
		parts := strings.Split(data, "_")
		if len(parts) >= 4 && parts[2] == "ex" {
			seconds, _ := strconv.Atoi(parts[3])
			exerciseID, _ := strconv.ParseInt(parts[5], 10, 64)
			startRestTimerWithExercise(chatID, seconds, exerciseID)
		}

	case strings.HasPrefix(data, "custom_timer_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "custom_timer_ex_"), 10, 64)
		userStates[chatID] = fmt.Sprintf("awaiting_custom_timer_ex_%d", exerciseID)
		msg := tgbotapi.NewMessage(chatID, "Введите время отдыха в секундах:")
		bot.Send(msg)

	case strings.HasPrefix(data, "stats_"):
		period := strings.TrimPrefix(data, "stats_")
		showStatistics(chatID, period)

	case data == "back_to_menu":
		sendMainMenu(chatID)

	case data == "my_workouts" || data == "create_new_workout":
		showMyWorkouts(chatID)

	case data == "setting_rest_timer":
		showRestTimerSettings(chatID)

	case data == "setting_weight_unit":
		showWeightUnitSettings(chatID)

	case data == "setting_notifications":
		showNotificationSettings(chatID)

	case data == "setting_export":
		showExportOptions(chatID)

	case data == "settings_back":
		showSettingsMenu(chatID)

	case strings.HasPrefix(data, "set_timer_"):
		secondsStr := strings.TrimPrefix(data, "set_timer_")
		seconds, _ := strconv.Atoi(secondsStr)
		saveUserSetting(chatID, "rest_timer", secondsStr)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Таймер отдыха установлен на %d секунд", seconds))
		bot.Send(msg)

	case data == "custom_timer":
		userStates[chatID] = "awaiting_custom_timer"
		msg := tgbotapi.NewMessage(chatID, "Введите время отдыха в секундах (например: 75):")
		bot.Send(msg)

	case strings.HasPrefix(data, "set_unit_"):
		unit := strings.TrimPrefix(data, "set_unit_")
		saveUserSetting(chatID, "weight_unit", unit)
		unitName := "килограммы"
		if unit == "lb" {
			unitName = "фунты"
		}
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Единицы измерения установлены: %s", unitName))
		bot.Send(msg)

	case data == "notifications_on":
		saveUserSetting(chatID, "notifications", "on")
		msg := tgbotapi.NewMessage(chatID, "✅ Уведомления включены")
		bot.Send(msg)

	case data == "notifications_off":
		saveUserSetting(chatID, "notifications", "off")
		msg := tgbotapi.NewMessage(chatID, "❌ Уведомления выключены")
		bot.Send(msg)

	case data == "notifications_time":
		userStates[chatID] = "awaiting_notification_time"
		msg := tgbotapi.NewMessage(chatID, "Введите время уведомлений (формат: 09:00 или 18:30):")
		bot.Send(msg)

	case strings.HasPrefix(data, "export_"):
		format := strings.TrimPrefix(data, "export_")
		startExport(chatID, format)
	}

	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	bot.Request(callbackConfig)
}

// ================== Основные функции интерфейса ==================
func sendMainMenu(chatID int64) {
	text := "🏋️‍♂️ *Добро пожаловать в Workout Bot!*\n\nВыберите действие:"

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Создать тренировку"),
			tgbotapi.NewKeyboardButton("▶️ Начать тренировку"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Мои тренировки"),
			tgbotapi.NewKeyboardButton("📊 Статистика"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⚙️ Настройки"),
		),
	)
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showWorkoutTypeMenu(chatID int64) {
	text := "Выберите тип тренировки:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🦵 Ноги", "create_workout_legs"),
			tgbotapi.NewInlineKeyboardButtonData("💪 Руки", "create_workout_arms"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Спина", "create_workout_back"),
			tgbotapi.NewInlineKeyboardButtonData("🫀 Грудь", "create_workout_chest"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌀 Плечи", "create_workout_shoulders"),
			tgbotapi.NewInlineKeyboardButtonData("⚡️ Кардио", "create_workout_cardio"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func createWorkoutDay(chatID int64, workoutType string) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	workoutDay := WorkoutDay{
		UserID:    user.ID,
		Name:      workoutType,
		StartedAt: time.Now(),
		Completed: false,
	}

	switch workoutType {
	case "legs":
		workoutDay.Exercises = getLegExercises()
	case "arms":
		workoutDay.Exercises = getArmExercises()
	case "back":
		workoutDay.Exercises = getBackExercises()
	case "chest":
		workoutDay.Exercises = getChestExercises()
	case "shoulders":
		workoutDay.Exercises = getShoulderExercises()
	case "cardio":
		workoutDay.Exercises = getCardioExercises()
	default:
		workoutDay.Exercises = getDefaultExercises()
	}

	db.Create(&workoutDay)
	showCreatedWorkout(chatID, workoutDay.ID)
}

func showCreatedWorkout(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)
	log.Println("workout: %v", workoutDay)

	var exercisesText strings.Builder
	exercisesText.WriteString(fmt.Sprintf("✅ *Тренировка создана: %s*\n\n", workoutDay.Name))
	exercisesText.WriteString("*Упражнения:*\n")

	for i, exercise := range workoutDay.Exercises {
		exercisesText.WriteString(fmt.Sprintf("%d. %s - %d подходов × %d повторений\n",
			i+1, exercise.Name, exercise.TargetSets, exercise.TargetReps))
	}

	exercisesText.WriteString("\nВыберите действие:")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Начать тренировку", fmt.Sprintf("start_workout_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать", fmt.Sprintf("edit_workout_%d", workoutDay.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить упражнение", fmt.Sprintf("add_exercise_to_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("delete_workout_%d", workoutDay.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, exercisesText.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showCurrentExerciseSession(chatID int64, workoutDayID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutDayID)

	if len(workoutDay.Exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
		bot.Send(msg)
		return
	}

	var session WorkoutSession
	db.Where("workout_day_id = ? AND is_active = ?", workoutDayID, true).
		Order("created_at DESC").
		First(&session)

	exerciseIndex := session.CurrentExerciseIndex
	if exerciseIndex >= len(workoutDay.Exercises) {
		exerciseIndex = 0
	}

	exercise := workoutDay.Exercises[exerciseIndex]

	var completedSets int64
	db.Model(&Set{}).Where("exercise_id = ? AND completed = ?", exercise.ID, true).Count(&completedSets)

	text := fmt.Sprintf(
		"🏋️‍♂️ *Тренировка: %s*\n\n"+
			"*Упражнение %d/%d:* %s\n\n"+
			"Цель: %d подходов × %d повторений\n"+
			"Выполнено: %d/%d подходов\n\n"+
			"*Что делаем?*",
		workoutDay.Name,
		exerciseIndex+1, len(workoutDay.Exercises), exercise.Name,
		exercise.TargetSets, exercise.TargetReps,
		completedSets, exercise.TargetSets,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Завершить подход",
				fmt.Sprintf("complete_set_ex_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Больше повторений",
				fmt.Sprintf("add_reps_ex_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("⚖️ Изменить вес",
				fmt.Sprintf("change_weight_ex_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏸️ Таймер отдыха",
				fmt.Sprintf("rest_timer_ex_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("➡️ След. упр-е",
				fmt.Sprintf("next_exercise_%d", workoutDayID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Прогресс",
				fmt.Sprintf("show_progress_%d", workoutDayID)),
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
				fmt.Sprintf("finish_workout_id_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func updateExerciseProgress(chatID int64, exerciseID int64) {
	var exercise Exercise
	db.Preload("Sets").First(&exercise, exerciseID)

	completedSets := 0
	for _, set := range exercise.Sets {
		if set.Completed {
			completedSets++
		}
	}

	text := fmt.Sprintf(
		"*%s*\n\nЦель: %d подходов по %d повторений\n\nВыполнено подходов: %d/%d\n\nСледующий подход через: 90 сек",
		exercise.Name, exercise.TargetSets, exercise.TargetReps, completedSets, exercise.TargetSets,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Завершить подход", fmt.Sprintf("complete_set_%d_%d", exercise.ID, completedSets)),
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить повтор", fmt.Sprintf("add_reps_%d", exercise.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏸️ Пауза (60 сек)", "pause_60"),
			tgbotapi.NewInlineKeyboardButtonData("⏹️ Завершить упражнение", "finish_exercise"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showStatsMenu(chatID int64, userID int64) {
	text := "📊 *Статистика тренировок*\n\nВыберите период:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📅 За неделю", "stats_week"),
			tgbotapi.NewInlineKeyboardButtonData("🗓️ За месяц", "stats_month"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📈 Общая статистика", "stats_all"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showStatistics(chatID int64, period string) {
	var statsText string

	switch period {
	case "week":
		statsText = "📅 *Статистика за неделю*\n\n✅ Тренировок: 3\n🏋️‍♂️ Упражнений: 15\n🔥 Подходов: 45\n⏱️ Среднее время: 45 мин"
	case "month":
		statsText = "🗓️ *Статистика за месяц*\n\n✅ Тренировок: 12\n🏋️‍♂️ Упражнений: 60\n🔥 Подходов: 180\n📈 Прогресс: +15%"
	default:
		statsText = "📈 *Общая статистика*\n\nВсего тренировок: 45\nРекорд веса: 120 кг\nЛюбимая группа: Ноги"
	}

	msg := tgbotapi.NewMessage(chatID, statsText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// ================== Вспомогательные функции ==================
func getUser(chatID int64, username string) (*User, error) {
	var user User
	result := db.Where("chat_id = ?", chatID).First(&user)

	if result.Error != nil {
		user = User{
			ChatID:    chatID,
			Username:  username,
			CreatedAt: time.Now(),
		}
		db.Create(&user)
		log.Println("created user")
	} else {
		log.Println("found user")
	}

	return &user, nil
}

func addRepsToLastSet(chatID int64, exerciseID int64) {
	var lastSet Set
	db.Where("exercise_id = ? AND completed = ?", exerciseID, true).
		Order("completed_at DESC").
		First(&lastSet)

	if lastSet.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Нет завершенных подходов для этого упражнения.")
		bot.Send(msg)
		return
	}

	lastSet.Reps += 1
	db.Save(&lastSet)

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"✅ Добавлено повторение!\n\nПодход №%d: %d повторений (вес: %.1f кг)",
		lastSet.ID, lastSet.Reps, lastSet.Weight,
	))
	bot.Send(msg)
	updateExerciseProgress(chatID, exerciseID)
}

func startExercise(chatID int64, exerciseID int64) {
	var exercise Exercise
	db.First(&exercise, exerciseID)

	var sets []Set
	for i := 0; i < exercise.TargetSets; i++ {
		sets = append(sets, Set{
			ExerciseID: exerciseID,
			Reps:       exercise.TargetReps,
			Weight:     0,
			Completed:  false,
		})
	}
	db.Create(&sets)

	text := fmt.Sprintf(
		"🎯 *Начинаем упражнение: %s*\n\n"+
			"Цель: %d подходов × %d повторений\n\n"+
			"Нажмите '✅ Завершить подход', когда выполните подход.",
		exercise.Name, exercise.TargetSets, exercise.TargetReps,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"✅ Завершить подход 1",
				fmt.Sprintf("complete_set_%d_%d", exerciseID, 1),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏸️ Таймер отдыха (90с)", "timer_90"),
			tgbotapi.NewInlineKeyboardButtonData("⚖️ Изменить вес", fmt.Sprintf("change_weight_%d", exerciseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Сменить упражнение", "change_exercise"),
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить тренировку", "finish_workout"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func finishWorkout(chatID int64) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	var workoutDay WorkoutDay
	db.Where("user_id = ? AND completed = ?", user.ID, false).
		Order("created_at DESC").
		First(&workoutDay)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ У вас нет активных тренировок.")
		bot.Send(msg)
		return
	}

	now := time.Now()
	workoutDay.Completed = true
	workoutDay.EndedAt = &now
	db.Save(&workoutDay)

	var exercises []Exercise
	db.Where("workout_day_id = ?", workoutDay.ID).Find(&exercises)

	totalSets := 0
	totalReps := 0
	completedExercises := 0

	for _, exercise := range exercises {
		var sets []Set
		db.Where("exercise_id = ?", exercise.ID).Find(&sets)

		completedSets := 0
		for _, set := range sets {
			if set.Completed {
				completedSets++
				totalReps += set.Reps
			}
		}
		totalSets += completedSets
		if completedSets > 0 {
			completedExercises++
		}
	}

	duration := now.Sub(workoutDay.StartedAt)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	text := fmt.Sprintf(
		"🏁 *Тренировка завершена!*\n\n"+
			"📊 *Итоги:*\n"+
			"• Время: %dч %dмин\n"+
			"• Упражнений: %d/%d\n"+
			"• Подходов: %d\n"+
			"• Повторений: %d\n\n"+
			"💪 Отличная работа! Отдыхайте и восстанавливайтесь!",
		hours, minutes, completedExercises, len(exercises), totalSets, totalReps,
	)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Посмотреть статистику"),
			tgbotapi.NewKeyboardButton("➕ Новая тренировка"),
		),
	)
	keyboard.ResizeKeyboard = true

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func startRestTimer(chatID int64, seconds int) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\nОтдыхайте!", seconds))
	message, _ := bot.Send(msg)

	ticker := time.NewTicker(10 * time.Second)
	timeUp := time.After(time.Duration(seconds) * time.Second)
	remaining := seconds

	go func() {
		for {
			select {
			case <-ticker.C:
				remaining -= 10
				if remaining > 0 {
					editMsg := tgbotapi.NewEditMessageText(
						chatID,
						message.MessageID,
						fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\nОтдыхайте!", remaining),
					)
					bot.Send(editMsg)
				}
			case <-timeUp:
				ticker.Stop()
				editMsg := tgbotapi.NewEditMessageText(
					chatID,
					message.MessageID,
					"🔔 *Время отдыха закончилось!*\n\nПриступайте к следующему подходу!",
				)
				editMsg.ParseMode = "Markdown"
				bot.Send(editMsg)
				return
			}
		}
	}()
}

func startRestTimerWithExercise(chatID int64, seconds int, exerciseID int64) {
	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\nРасслабьтесь и подготовьтесь к следующему подходу!", seconds))

	message, _ := bot.Send(msg)

	go func() {
		remaining := seconds

		for remaining > 0 {
			time.Sleep(1 * time.Second)
			remaining--

			if remaining%10 == 0 || remaining <= 5 {
				editMsg := tgbotapi.NewEditMessageText(
					chatID,
					message.MessageID,
					fmt.Sprintf("⏳ Таймер отдыха: %d секунд\n\nРасслабьтесь и подготовьтесь к следующему подходу!", remaining),
				)
				bot.Send(editMsg)
			}
		}

		editMsg := tgbotapi.NewEditMessageText(
			chatID,
			message.MessageID,
			"🔔 *Время отдыха закончилось!*\n\nПриступайте к следующему подходу! 💪",
		)
		editMsg.ParseMode = "Markdown"

		editMarkup := tgbotapi.NewEditMessageReplyMarkup(
			chatID,
			message.MessageID,
			tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("✅ Начать подход",
						fmt.Sprintf("complete_set_ex_%d", exerciseID)),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("➕ Повторения",
						fmt.Sprintf("add_reps_ex_%d", exerciseID)),
					tgbotapi.NewInlineKeyboardButtonData("⚖️ Вес",
						fmt.Sprintf("change_weight_ex_%d", exerciseID)),
				),
			),
		)

		bot.Send(editMsg)
		bot.Send(editMarkup)
	}()
}

func askForNewWeight(chatID int64, exerciseID int64) {
	userStates[chatID] = fmt.Sprintf("awaiting_weight_%d", exerciseID)
	msg := tgbotapi.NewMessage(chatID, "⚖️ Введите новый вес (в кг):")
	bot.Send(msg)
}

func showSettingsMenu(chatID int64) {
	text := "⚙️ *Настройки*\n\nВыберите параметр для изменения:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏱️ Таймер отдыха", "setting_rest_timer"),
			tgbotapi.NewInlineKeyboardButtonData("🏋️‍♂️ Единицы веса", "setting_weight_unit"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔔 Уведомления", "setting_notifications"),
			tgbotapi.NewInlineKeyboardButtonData("📁 Экспорт данных", "setting_export"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад в меню", "back_to_menu"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showRestTimerSettings(chatID int64) {
	text := "⏱️ *Настройка таймера отдыха*\n\nВыберите продолжительность отдыха между подходами:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("60 сек", "set_timer_60"),
			tgbotapi.NewInlineKeyboardButtonData("90 сек", "set_timer_90"),
			tgbotapi.NewInlineKeyboardButtonData("120 сек", "set_timer_120"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("180 сек", "set_timer_180"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Свое время", "custom_timer"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings_back"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showWeightUnitSettings(chatID int64) {
	text := "🏋️‍♂️ *Единицы измерения веса*\n\nВыберите предпочитаемые единицы:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("кг (килограммы)", "set_unit_kg"),
			tgbotapi.NewInlineKeyboardButtonData("lb (фунты)", "set_unit_lb"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings_back"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showNotificationSettings(chatID int64) {
	text := "🔔 *Настройка уведомлений*\n\nУправляйте напоминаниями о тренировках:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Включить", "notifications_on"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Выключить", "notifications_off"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏰ Настроить время", "notifications_time"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings_back"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showExportOptions(chatID int64) {
	text := "📁 *Экспорт данных*\n\nВыберите формат экспорта:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📄 CSV", "export_csv"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Excel", "export_excel"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Google Sheets", "export_gsheets"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "settings_back"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// ================== Функции для шаблонов упражнений ==================
func getLegExercises() []Exercise {
	return []Exercise{
		{Name: "Приседания со штангой", TargetSets: 4, TargetReps: 10},
		{Name: "Жим ногами", TargetSets: 3, TargetReps: 12},
		{Name: "Выпады", TargetSets: 3, TargetReps: 12},
		{Name: "Сгибания ног", TargetSets: 3, TargetReps: 15},
		{Name: "Разгибания ног", TargetSets: 3, TargetReps: 15},
	}
}

func getArmExercises() []Exercise {
	return []Exercise{
		{Name: "Подъем штанги на бицепс", TargetSets: 4, TargetReps: 10},
		{Name: "Молотки", TargetSets: 3, TargetReps: 12},
		{Name: "Концентрированные сгибания", TargetSets: 3, TargetReps: 12},
		{Name: "Отжимания на брусьях", TargetSets: 4, TargetReps: 10},
		{Name: "Французский жим", TargetSets: 3, TargetReps: 12},
	}
}

func getBackExercises() []Exercise {
	return []Exercise{
		{Name: "Становая тяга", TargetSets: 4, TargetReps: 8},
		{Name: "Тяга штанги в наклоне", TargetSets: 4, TargetReps: 10},
		{Name: "Подтягивания", TargetSets: 4, TargetReps: 10},
		{Name: "Тяга верхнего блока", TargetSets: 3, TargetReps: 12},
		{Name: "Гиперэкстензия", TargetSets: 3, TargetReps: 15},
	}
}

func getChestExercises() []Exercise {
	return []Exercise{
		{Name: "Жим штанги лежа", TargetSets: 4, TargetReps: 10},
		{Name: "Жим гантелей", TargetSets: 3, TargetReps: 12},
		{Name: "Разводка гантелей", TargetSets: 3, TargetReps: 12},
		{Name: "Отжимания", TargetSets: 4, TargetReps: 15},
		{Name: "Сведения в кроссовере", TargetSets: 3, TargetReps: 15},
	}
}

func getShoulderExercises() []Exercise {
	return []Exercise{
		{Name: "Жим штанги с груди", TargetSets: 4, TargetReps: 10},
		{Name: "Махи гантелями в стороны", TargetSets: 3, TargetReps: 12},
		{Name: "Махи в наклоне", TargetSets: 3, TargetReps: 12},
		{Name: "Тяга штанги к подбородку", TargetSets: 3, TargetReps: 12},
		{Name: "Подъемы гантелей перед собой", TargetSets: 3, TargetReps: 12},
	}
}

func getCardioExercises() []Exercise {
	return []Exercise{
		{Name: "Беговая дорожка", TargetSets: 1, TargetReps: 20},
		{Name: "Велотренажер", TargetSets: 1, TargetReps: 20},
		{Name: "Скакалка", TargetSets: 5, TargetReps: 100},
		{Name: "Гребной тренажер", TargetSets: 3, TargetReps: 10},
	}
}

func getDefaultExercises() []Exercise {
	return []Exercise{
		{Name: "Базовое упражнение 1", TargetSets: 3, TargetReps: 10},
		{Name: "Базовое упражнение 2", TargetSets: 3, TargetReps: 10},
		{Name: "Базовое упражнение 3", TargetSets: 3, TargetReps: 10},
	}
}

// ================== Функции управления тренировками ==================
func showMyWorkouts(chatID int64) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)
	log.Println("found user: %v", user)

	var workouts []WorkoutDay
	db.Where("user_id = ?", user.ID).Order("created_at DESC").Find(&workouts)
	log.Println("found workouts: %v", workouts)

	if len(workouts) == 0 {
		msg := tgbotapi.NewMessage(chatID, "📭 У вас пока нет созданных тренировок.\n\nСоздайте первую тренировку!")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать тренировку", "create_new_workout"),
			),
		)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		return
	}

	text := "📋 *Ваши тренировки:*\n\n"
	for i, workout := range workouts {
		status := "🟢 Активна"
		if workout.Completed {
			status = "✅ Завершена"
		}
		date := workout.StartedAt.Format("02.01.2006")
		text += fmt.Sprintf("%d. *%s* - %s\n   📅 %s\n\n",
			i+1, workout.Name, status, date)
	}

	text += "Выберите тренировку для просмотра:"

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, workout := range workouts {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rowIndex := len(rows) - 1
		buttonText := fmt.Sprintf("%s %d", workout.Name, i+1)
		rows[rowIndex] = append(rows[rowIndex],
			tgbotapi.NewInlineKeyboardButtonData(buttonText,
				fmt.Sprintf("view_workout_%d", workout.ID)))
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("➕ Новая тренировка", "create_new_workout"),
		tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
	})

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showWorkoutDetails(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		bot.Send(msg)
		return
	}

	var text strings.Builder
	status := "🟢 Активна"
	if workoutDay.Completed {
		status = "✅ Завершена"
		endDate := ""
		if workoutDay.EndedAt != nil {
			endDate = workoutDay.EndedAt.Format("15:04")
		}
		text.WriteString(fmt.Sprintf("✅ *%s* (Завершена в %s)\n\n", workoutDay.Name, endDate))
	} else {
		text.WriteString(fmt.Sprintf("🟢 *%s* (Активна)\n\n", workoutDay.Name))
	}

	text.WriteString(fmt.Sprintf("Статус: %s\n", status))
	text.WriteString(fmt.Sprintf("Дата: %s\n\n", workoutDay.StartedAt.Format("02.01.2006")))

	text.WriteString("*Упражнения:*\n")
	for i, exercise := range workoutDay.Exercises {
		var completedSets int64
		db.Model(&Set{}).Where("exercise_id = ? AND completed = ?", exercise.ID, true).Count(&completedSets)

		text.WriteString(fmt.Sprintf("%d. %s: %d/%d подходов\n",
			i+1, exercise.Name, completedSets, exercise.TargetSets))
	}

	var keyboard tgbotapi.InlineKeyboardMarkup

	if !workoutDay.Completed {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("▶️ Начать тренировку",
					fmt.Sprintf("start_active_workout_%d", workoutDay.ID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать",
					fmt.Sprintf("edit_workout_%d", workoutDay.ID)),
				tgbotapi.NewInlineKeyboardButtonData("➕ Упр-е",
					fmt.Sprintf("add_exercise_to_%d", workoutDay.ID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить",
					fmt.Sprintf("confirm_delete_%d", workoutDay.ID)),
				tgbotapi.NewInlineKeyboardButtonData("📋 Все тренировки", "my_workouts"),
			),
		)
	} else {
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📊 Статистика",
					fmt.Sprintf("stats_workout_%d", workoutDay.ID)),
				tgbotapi.NewInlineKeyboardButtonData("🔄 Повторить",
					fmt.Sprintf("repeat_workout_%d", workoutDay.ID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 Все тренировки", "my_workouts"),
				tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func startSpecificWorkout(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)

	if workoutDay.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
		bot.Send(msg)
		return
	}

	if workoutDay.Completed {
		msg := tgbotapi.NewMessage(chatID, "❌ Эта тренировка уже завершена. Создайте новую или повторите эту.")
		bot.Send(msg)
		return
	}

	session := WorkoutSession{
		WorkoutDayID:         workoutDay.ID,
		StartedAt:            time.Now(),
		IsActive:             true,
		CurrentExerciseIndex: 0,
	}
	db.Create(&session)
	showCurrentExerciseSession(chatID, workoutDay.ID)
}

func startActiveWorkout(chatID int64, userID int64) {
	var workouts []WorkoutDay
	db.Where("user_id = ? AND completed = ?", userID, false).
		Order("created_at DESC").
		Find(&workouts)

	if len(workouts) == 0 {
		msg := tgbotapi.NewMessage(chatID,
			"У вас нет активных тренировок. Сначала создайте тренировку!")
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➕ Создать тренировку", "create_new_workout"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
			),
		)
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		return
	}

	text := "▶️ *Выберите тренировку для начала:*\n\n"
	for i, workout := range workouts {
		text += fmt.Sprintf("%d. *%s* (создана %s)\n",
			i+1, workout.Name, workout.StartedAt.Format("02.01"))
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, workout := range workouts {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rowIndex := len(rows) - 1
		buttonText := fmt.Sprintf("%s", workout.Name)
		rows[rowIndex] = append(rows[rowIndex],
			tgbotapi.NewInlineKeyboardButtonData(buttonText,
				fmt.Sprintf("start_active_workout_%d", workout.ID)))
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("📋 Все тренировки", "my_workouts"),
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_menu"),
	})

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

// ================== Функции редактирования и управления ==================
func editWorkout(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)

	text := fmt.Sprintf("✏️ *Редактирование: %s*\n\n", workoutDay.Name)
	text += "Выберите что изменить:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Изменить название",
				fmt.Sprintf("edit_name_%d", workoutDay.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить упражнение",
				fmt.Sprintf("add_exercise_to_%d", workoutDay.ID)),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Упр-я",
				fmt.Sprintf("edit_exercises_%d", workoutDay.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад",
				fmt.Sprintf("view_workout_%d", workoutDay.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func askForNewExercise(chatID int64, workoutID int64) {
	userStates[chatID] = fmt.Sprintf("awaiting_exercise_name_%d", workoutID)

	text := "➕ *Добавление упражнения*\n\n" +
		"Введите название нового упражнения:"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Отмена",
				fmt.Sprintf("view_workout_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func confirmDeleteWorkout(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.First(&workoutDay, workoutID)

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
	bot.Send(msg)
}

func deleteWorkout(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)

	for _, exercise := range workoutDay.Exercises {
		db.Where("exercise_id = ?", exercise.ID).Delete(&Set{})
	}

	db.Where("workout_day_id = ?", workoutID).Delete(&Exercise{})
	db.Where("workout_day_id = ?", workoutID).Delete(&WorkoutSession{})
	db.Delete(&workoutDay)

	msg := tgbotapi.NewMessage(chatID, "✅ Тренировка успешно удалена!")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 В меню", "back_to_menu"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showWorkoutStatistics(chatID int64, workoutID int64) {
	var workoutDay WorkoutDay
	db.Preload("Exercises").First(&workoutDay, workoutID)

	var exercises []Exercise
	db.Where("workout_day_id = ?", workoutID).Find(&exercises)

	totalSets := 0
	totalReps := 0
	totalWeight := 0.0
	completedExercises := 0

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Статистика: %s*\n\n", workoutDay.Name))

	if workoutDay.EndedAt != nil {
		duration := workoutDay.EndedAt.Sub(workoutDay.StartedAt)
		text.WriteString(fmt.Sprintf("⏱️ *Время:* %s\n", formatDuration(duration)))
	}

	text.WriteString(fmt.Sprintf("📅 *Дата:* %s\n\n", workoutDay.StartedAt.Format("02.01.2006 15:04")))

	for _, exercise := range exercises {
		var sets []Set
		db.Where("exercise_id = ? AND completed = ?", exercise.ID, true).Find(&sets)

		if len(sets) == 0 {
			continue
		}

		completedExercises++
		exerciseReps := 0
		exerciseWeight := 0.0

		for _, set := range sets {
			exerciseReps += set.Reps
			exerciseWeight += float64(set.Weight)
			totalReps += set.Reps
			totalWeight += float64(set.Weight)
		}

		totalSets += len(sets)

		avgWeight := 0.0
		if len(sets) > 0 {
			avgWeight = exerciseWeight / float64(len(sets))
		}

		text.WriteString(fmt.Sprintf("• *%s:* %d×%d = %d повторений (avg %.1f кг)\n",
			exercise.Name, len(sets), exercise.TargetReps, exerciseReps, avgWeight))
	}

	text.WriteString(fmt.Sprintf("\n📈 *Итого:*\n"))
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, len(exercises)))
	text.WriteString(fmt.Sprintf("• Подходов: %d\n", totalSets))
	text.WriteString(fmt.Sprintf("• Повторений: %d\n", totalReps))

	if totalSets > 0 {
		avgTotalWeight := totalWeight / float64(totalSets)
		text.WriteString(fmt.Sprintf("• Средний вес: %.1f кг\n", avgTotalWeight))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Повторить тренировку",
				fmt.Sprintf("repeat_workout_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад",
				fmt.Sprintf("view_workout_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
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

func repeatWorkout(chatID int64, workoutID int64) {
	var originalWorkout WorkoutDay
	db.Preload("Exercises").First(&originalWorkout, workoutID)

	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	newWorkout := WorkoutDay{
		UserID:    user.ID,
		Name:      fmt.Sprintf("%s (повтор)", originalWorkout.Name),
		StartedAt: time.Now(),
		Completed: false,
	}

	var exercises []Exercise
	for _, originalExercise := range originalWorkout.Exercises {
		exercise := Exercise{
			Name:       originalExercise.Name,
			TargetSets: originalExercise.TargetSets,
			TargetReps: originalExercise.TargetReps,
		}
		exercises = append(exercises, exercise)
	}

	db.Create(&newWorkout)
	for i := range exercises {
		exercises[i].WorkoutDayID = newWorkout.ID
		db.Create(&exercises[i])
	}

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("🔄 Тренировка *%s* скопирована!\n\nВы можете начать её прямо сейчас.", newWorkout.Name))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Начать",
				fmt.Sprintf("start_workout_%d", newWorkout.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать",
				fmt.Sprintf("edit_workout_%d", newWorkout.ID)),
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои тренировки", "my_workouts"),
		),
	)

	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func completeExerciseSet(chatID int64, exerciseID int64) {
	var exercise Exercise
	db.First(&exercise, exerciseID)

	var nextSet Set
	db.Where("exercise_id = ? AND completed = ?", exerciseID, false).
		Order("created_at ASC").
		First(&nextSet)

	var set Set
	if nextSet.ID == 0 {
		var lastWeightSet Set
		db.Where("exercise_id = ? AND weight > ?", exerciseID, 0).
			Order("completed_at DESC").
			First(&lastWeightSet)

		weight := float32(0)
		if lastWeightSet.ID != 0 {
			weight = lastWeightSet.Weight
		}

		set = Set{
			ExerciseID:  exerciseID,
			Reps:        exercise.TargetReps,
			Weight:      weight,
			Completed:   true,
			CompletedAt: &[]time.Time{time.Now()}[0],
		}
		db.Create(&set)
	} else {
		set = nextSet
		set.Completed = true
		set.CompletedAt = &[]time.Time{time.Now()}[0]
		set.Reps = exercise.TargetReps
		db.Save(&set)
	}

	var completedSets int64
	db.Model(&Set{}).Where("exercise_id = ? AND completed = ?", exerciseID, true).Count(&completedSets)

	var workoutDay WorkoutDay
	db.Joins("JOIN exercises ON exercises.workout_day_id = workout_days.id").
		Where("exercises.id = ?", exerciseID).
		First(&workoutDay)

	text := fmt.Sprintf("✅ *Подход завершен!*\n\n"+
		"Упражнение: %s\n"+
		"Подход: %d/%d\n"+
		"Повторений: %d\n"+
		"Вес: %.1f кг\n\n"+
		"Отдыхайте %d секунд перед следующим подходом.",
		exercise.Name, completedSets, exercise.TargetSets,
		set.Reps, set.Weight, 90)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏱️ Таймер 90с",
				fmt.Sprintf("timer_90_ex_%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("➕ Повторения",
				fmt.Sprintf("add_reps_ex_%d", exerciseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚖️ Изменить вес",
				fmt.Sprintf("change_weight_ex_%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("➡️ Следующее",
				fmt.Sprintf("next_exercise_wd_%d", workoutDay.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 К упражнению",
				fmt.Sprintf("show_exercise_%d", exerciseID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func showRestTimerOptions(chatID int64, exerciseID int64) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	var setting UserSetting
	db.Where("user_id = ? AND key = ?", user.ID, "rest_timer").First(&setting)

	defaultTimer := 90
	if setting.ID != 0 {
		if seconds, err := strconv.Atoi(setting.Value); err == nil {
			defaultTimer = seconds
		}
	}

	text := fmt.Sprintf("⏱️ *Таймер отдыха*\n\n"+
		"Текущая настройка: %d секунд\n\n"+
		"Выберите время отдыха:", defaultTimer)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("60 сек",
				fmt.Sprintf("start_timer_60_ex_%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("90 сек",
				fmt.Sprintf("start_timer_90_ex_%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("120 сек",
				fmt.Sprintf("start_timer_120_ex_%d", exerciseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("180 сек",
				fmt.Sprintf("start_timer_180_ex_%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Свое",
				fmt.Sprintf("custom_timer_ex_%d", exerciseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "setting_rest_timer"),
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад",
				fmt.Sprintf("show_exercise_%d", exerciseID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func moveToNextExercise(chatID int64, workoutDayID int64) {
	var session WorkoutSession
	db.Where("workout_day_id = ? AND is_active = ?", workoutDayID, true).
		Order("created_at DESC").
		First(&session)

	if session.ID == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Активная сессия не найдена")
		bot.Send(msg)
		return
	}

	var exercises []Exercise
	db.Where("workout_day_id = ?", workoutDayID).Find(&exercises)

	if len(exercises) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ В тренировке нет упражнений")
		bot.Send(msg)
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
					fmt.Sprintf("finish_workout_id_%d", workoutDayID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Начать заново",
					fmt.Sprintf("restart_workout_%d", workoutDayID)),
				tgbotapi.NewInlineKeyboardButtonData("🔙 К первому",
					fmt.Sprintf("first_exercise_%d", workoutDayID)),
			),
		)

		msg.ReplyMarkup = keyboard
		bot.Send(msg)
		return
	}

	db.Save(&session)
	showCurrentExerciseSession(chatID, workoutDayID)
}

func showWorkoutProgress(chatID int64, workoutDayID int64) {
	var workoutDay WorkoutDay
	db.First(&workoutDay, workoutDayID)

	var exercises []Exercise
	db.Where("workout_day_id = ?", workoutDayID).Find(&exercises)

	var text strings.Builder
	text.WriteString(fmt.Sprintf("📊 *Прогресс тренировки: %s*\n\n", workoutDay.Name))

	totalExercises := len(exercises)
	completedExercises := 0
	totalSets := 0
	completedSets := 0

	for i, exercise := range exercises {
		var completedExerciseSets int64
		db.Model(&Set{}).Where("exercise_id = ? AND completed = ?", exercise.ID, true).Count(&completedExerciseSets)

		var allSets int64
		db.Model(&Set{}).Where("exercise_id = ?", exercise.ID).Count(&allSets)

		if allSets == 0 {
			allSets = int64(exercise.TargetSets)
		}

		status := "🔴"
		if int(completedExerciseSets) >= exercise.TargetSets {
			status = "✅"
			completedExercises++
		} else if completedExerciseSets > 0 {
			status = "🟡"
		}

		text.WriteString(fmt.Sprintf("%s %d. %s: %d/%d подходов\n",
			status, i+1, exercise.Name, completedExerciseSets, exercise.TargetSets))

		completedSets += int(completedExerciseSets)
		totalSets += exercise.TargetSets
	}

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
				fmt.Sprintf("continue_workout_%d", workoutDayID)),
			tgbotapi.NewInlineKeyboardButtonData("📊 Детали",
				fmt.Sprintf("detailed_stats_%d", workoutDayID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 К тренировке",
				fmt.Sprintf("view_workout_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func confirmFinishWorkout(chatID int64, workoutDayID int64) {
	var workoutDay WorkoutDay
	db.First(&workoutDay, workoutDayID)

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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Сначала статистика",
				fmt.Sprintf("pre_finish_stats_%d", workoutDayID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func finishWorkoutById(chatID int64, workoutDayID int64) {
	var workoutDay WorkoutDay
	db.First(&workoutDay, workoutDayID)

	now := time.Now()
	workoutDay.Completed = true
	workoutDay.EndedAt = &now
	db.Save(&workoutDay)

	db.Model(&WorkoutSession{}).
		Where("workout_day_id = ? AND is_active = ?", workoutDayID, true).
		Update("is_active", false)

	showWorkoutStatistics(chatID, workoutDayID)
}

func saveUserSetting(chatID int64, key string, value string) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	var setting UserSetting
	db.Where("user_id = ? AND key = ?", user.ID, key).First(&setting)

	if setting.ID == 0 {
		setting = UserSetting{
			UserID:    user.ID,
			Key:       key,
			Value:     value,
			UpdatedAt: time.Now(),
		}
		db.Create(&setting)
	} else {
		setting.Value = value
		setting.UpdatedAt = time.Now()
		db.Save(&setting)
	}
}

func startExport(chatID int64, format string) {
	var user User
	db.Where("chat_id = ?", chatID).First(&user)

	var workouts []WorkoutDay
	db.Where("user_id = ?", user.ID).Find(&workouts)

	msg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("📦 *Экспорт данных*\n\n"+
			"Формат: %s\n"+
			"Найдено тренировок: %d\n\n"+
			"Функция экспорта в разработке...",
			strings.ToUpper(format), len(workouts)))
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func handleState(chatID int64, userID int64, text string) {
	state, exists := userStates[chatID]
	if !exists {
		return
	}

	switch {
	case strings.HasPrefix(state, "awaiting_weight_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 3 {
			exerciseID, _ := strconv.ParseInt(parts[2], 10, 64)

			weight, err := strconv.ParseFloat(text, 32)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат веса. Введите число (например: 42.5)")
				bot.Send(msg)
				return
			}

			var lastSet Set
			db.Where("exercise_id = ?", exerciseID).
				Order("created_at DESC").
				First(&lastSet)

			if lastSet.ID != 0 {
				lastSet.Weight = float32(weight)
				db.Save(&lastSet)

				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
					"✅ Вес обновлен: %.1f кг для подхода №%d",
					weight, lastSet.ID,
				))
				bot.Send(msg)
			}

			userStates[chatID] = ""
		}

	case strings.HasPrefix(state, "awaiting_exercise_name_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 4 {
			workoutID, _ := strconv.ParseInt(parts[3], 10, 64)

			exercise := Exercise{
				WorkoutDayID: workoutID,
				Name:         text,
				TargetSets:   3,
				TargetReps:   10,
			}
			db.Create(&exercise)

			userStates[chatID] = ""
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Упражнение '%s' добавлено!", text))
			bot.Send(msg)
			showWorkoutDetails(chatID, workoutID)
		}

	case state == "awaiting_custom_timer":
		seconds, err := strconv.Atoi(text)
		if err != nil || seconds < 10 || seconds > 300 {
			msg := tgbotapi.NewMessage(chatID, "❌ Неверное время. Введите число от 10 до 300 секунд.")
			bot.Send(msg)
			return
		}

		saveUserSetting(chatID, "rest_timer", text)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Таймер отдыха установлен на %d секунд", seconds))
		bot.Send(msg)
		userStates[chatID] = ""

	case strings.HasPrefix(state, "awaiting_custom_timer_ex_"):
		parts := strings.Split(state, "_")
		if len(parts) >= 5 {
			exerciseID, _ := strconv.ParseInt(parts[4], 10, 64)

			seconds, err := strconv.Atoi(text)
			if err != nil || seconds < 10 || seconds > 300 {
				msg := tgbotapi.NewMessage(chatID, "❌ Неверное время. Введите число от 10 до 300 секунд.")
				bot.Send(msg)
				return
			}

			startRestTimerWithExercise(chatID, seconds, exerciseID)
			userStates[chatID] = ""
		}

	case state == "awaiting_notification_time":
		if !strings.Contains(text, ":") || len(text) != 5 {
			msg := tgbotapi.NewMessage(chatID, "❌ Неверный формат. Используйте ЧЧ:ММ (например: 09:00)")
			bot.Send(msg)
			return
		}

		saveUserSetting(chatID, "notification_time", text)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("✅ Время уведомлений установлено на %s", text))
		bot.Send(msg)
		userStates[chatID] = ""
	}
}
