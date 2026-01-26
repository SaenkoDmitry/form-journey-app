package presenter

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowCurrentSession(chatID int64, res *dto.CurrentExerciseSession) {
	var text strings.Builder

	dayType := res.DayType
	exerciseIndex := res.ExerciseIndex
	workoutDay := res.WorkoutDay
	exercise := res.Exercise
	exerciseObj := res.ExerciseObj
	workoutID := workoutDay.ID

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
			tgbotapi.NewInlineKeyboardButtonData(messages.DoneSet, fmt.Sprintf("set_complete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.AddSet, fmt.Sprintf("set_add_one_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.RemoveSet, fmt.Sprintf("set_remove_last_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Timer, fmt.Sprintf("timer_start_%d_ex_%d", exercise.ExerciseType.RestInSeconds, exercise.ID)),
		),
		changeSettingsButtons,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Technique, fmt.Sprintf("exercise_show_hint_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.EndWorkout, fmt.Sprintf("workout_confirm_finish_%d", workoutID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.Prev, fmt.Sprintf("exercise_move_to_prev_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Progress, fmt.Sprintf("workout_show_progress_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.DropExercise, fmt.Sprintf("exercise_confirm_delete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.Next, fmt.Sprintf("exercise_move_to_next_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard

	p.bot.Send(msg)
}

func (p *Presenter) ShowNoExercises(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ В этой тренировке нет упражнений.")
	p.bot.Send(msg)
}

func (p *Presenter) ShowNotFoundExercise(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ Упражнение не найдено.")
	p.bot.Send(msg)
}

func (p *Presenter) ShowSelectExerciseForProgramDayDialog(chatID, dayTypeID int64, group *dto.Group, exerciseTypes []models.ExerciseType) {
	text := fmt.Sprintf("*Тип: %s \n\n Выберите упражнение из списка:*", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("change_program_day_add_exercise_%d_%d", dayTypeID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowConfirmDeleteDialog(chatID int64, res *dto.ConfirmDeleteExercise) {
	exercise := res.Exercise
	exerciseObj := res.ExerciseObj
	text := fmt.Sprintf("🗑️ <b>Удаление упражнения из тренировочного дня</b>\n\n"+
		"Вы уверены, что хотите удалить упражнение:\n"+
		"<b>%s</b>?\n\n"+
		"⚠️ Это действие нельзя отменить!", exerciseObj.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("exercise_delete_%d", exercise.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("workout_start_%d", exercise.WorkoutDayID)),
		),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) CompleteAllExercises(chatID, workoutID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"🎉 Вы завершили все упражнения в этой тренировке!\n\n"+
			"Хотите завершить тренировку или добавить еще упражнения?")

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏁 Завершить",
				fmt.Sprintf("workout_confirm_finish_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("➕ Еще упражнение",
				fmt.Sprintf("exercise_add_for_current_workout_%d", workoutID)),
		),
	)
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowHint(chatID int64, res *dto.GetExercise) {
	exercise := res.Exercise
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", fmt.Sprintf("exercise_show_current_session_%d", exercise.WorkoutDayID)),
	))
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, utils.WrapYandexLink(exercise.ExerciseType.Url))
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) AddExerciseDialog(chatID, workoutID int64, groups []models.ExerciseGroupType) {
	text := messages.SelectGroupOfMuscle
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1], tgbotapi.NewInlineKeyboardButtonData(group.Name,
			fmt.Sprintf("exercise_select_for_current_workout_%d_%s", workoutID, group.Code)))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowSelectExerciseForCurrentWorkoutDialog(chatID, workoutID int64, group *dto.Group, exerciseTypes []models.ExerciseType) {
	text := fmt.Sprintf("<b>Тип:</b> %s \n\n <b>Выберите упражнение из списка:</b>", group.Name)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, exercise := range exerciseTypes {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				exercise.Name,
				fmt.Sprintf("exercise_add_specific_for_current_workout_%d_%d", workoutID, exercise.ID),
			),
		))
	}
	fmt.Println("rows", len(rows), rows)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	_, err := p.bot.Send(msg)
	fmt.Println("err:", err)
}
