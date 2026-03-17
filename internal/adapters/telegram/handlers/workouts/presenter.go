package workouts

import (
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"time"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowWorkoutProgress(chatID int64, progress *dto.WorkoutProgress, stats *dto.WorkoutStatistic, needShowButtons bool) {
	totalWeight := stats.TotalWeight
	totalTime := stats.CardioTime

	var text strings.Builder

	w := progress.Workout
	text.WriteString(fmt.Sprintf("<b>День:</b> <u>%s</u> \n", progress.Workout.DayTypeName))
	text.WriteString(fmt.Sprintf("<b>Начата:</b> %s\n", w.StartedAt))
	text.WriteString(fmt.Sprintf("<b>Статус:</b> %s\n", w.Status))
	if w.Completed {
		text.WriteString(fmt.Sprintf("<b>Длительность:</b> %s\n", w.Duration))
	}
	text.WriteString("\n")

	if len(w.Exercises) > 0 {
		text.WriteString("<b>УПРАЖНЕНИЯ:</b>\n")
	}

	for i, ex := range w.Exercises {
		text.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, ex.Name))

		for _, set := range ex.Sets {
			text.WriteString(set.FormattedString)
		}
		if ex.SumWeight > 0 {
			text.WriteString(fmt.Sprintf("<u>Общий вес</u>: %.0f кг\n", ex.SumWeight))
		}
		text.WriteString("\n")
	}

	text.WriteString("\n📈 <b>ПРОГРЕСС:</b>\n")
	text.WriteString(fmt.Sprintf(
		"• Упражнений: %d/%d\n",
		progress.CompletedExercises,
		progress.TotalExercises,
	))
	text.WriteString(fmt.Sprintf(
		"• Подходов: %d/%d\n",
		progress.CompletedSets,
		progress.TotalSets,
	))
	text.WriteString(fmt.Sprintf(
		"• Прогресс: %d%%\n",
		progress.ProgressPercent,
	))

	if totalWeight > 0 {
		text.WriteString(fmt.Sprintf("• Общий тоннаж: %.0f кг\n", totalWeight))
	}
	if totalTime > 0 {
		text.WriteString(fmt.Sprintf("• Время кардио: %d минут\n", totalTime))
	}

	text.WriteString(fmt.Sprintf("• [%s]\n\n", progressBar(progress.ProgressPercent)))

	if progress.RemainingMin != nil {
		text.WriteString(fmt.Sprintf(
			"⏰ <b>Прогноз окончания:</b> ~%d минут\n",
			*progress.RemainingMin,
		))
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	if needShowButtons {
		keyboard := p.buildKeyboard(progress)
		msg.ReplyMarkup = keyboard
	}

	_, _ = p.bot.Send(msg)
}

func (p *Presenter) WorkoutCreated(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "✅ <b>Тренировка создана!</b>\n\n")
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func progressBar(percent int) string {
	const barLength = 13

	filled := (percent * barLength) / 100
	var b strings.Builder

	for i := 0; i < barLength; i++ {
		if i < filled {
			b.WriteString("🏋️‍♂️")
		} else {
			b.WriteString("░")
		}
	}
	return b.String()
}

func (p *Presenter) buildKeyboard(data *dto.WorkoutProgress) tgbotapi.InlineKeyboardMarkup {
	workoutID := data.Workout.ID

	backTo := tgbotapi.NewInlineKeyboardButtonData(
		messages.BackTo,
		"workout_show_my",
	)

	deleteBtn := tgbotapi.NewInlineKeyboardButtonData(
		"🗑️ Удалить",
		fmt.Sprintf("workout_confirm_delete_%d", workoutID),
	)

	if !data.Workout.Completed {
		addExerciseBtn := tgbotapi.NewInlineKeyboardButtonData(
			"➕ Еще упражнение",
			fmt.Sprintf("exercise_add_for_current_workout_%d", workoutID),
		)

		toWorkoutBtn := tgbotapi.NewInlineKeyboardButtonData("▶️ Начать", fmt.Sprintf("workout_start_%d", workoutID))
		if data.SessionStarted {
			toWorkoutBtn = tgbotapi.NewInlineKeyboardButtonData("▶️ К тренировке", fmt.Sprintf("exercise_show_current_session_%d", workoutID))
		}

		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(addExerciseBtn, deleteBtn),
			tgbotapi.NewInlineKeyboardRow(toWorkoutBtn),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backTo, deleteBtn),
	)
}

func (p *Presenter) ShowNotFoundSpecific(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ Тренировка не найдена")
	p.bot.Send(msg)
}

func (p *Presenter) ShowAlreadyCompleted(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "❌ Эта тренировка уже завершена. Создайте новую или повторите эту.")
	p.bot.Send(msg)
}

func (p *Presenter) ShowNotFoundAll(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "📭 У вас пока нет созданных тренировок.\n\nСоздайте первую тренировку!")
	p.bot.Send(msg)
}

func (p *Presenter) ShowNotFoundAllForUser(chatID int64, user *models.User) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📭 У пользователя %s пока нет созданных тренировок.", user.ShortName()))
	p.bot.Send(msg)
}

func (p *Presenter) ShowConfirmDeleteWorkout(chatID int64, res *dto.ConfirmDeleteWorkout) {
	text := fmt.Sprintf("🗑️ *Удаление тренировки*\n\n"+
		"Вы уверены, что хотите удалить тренировку:\n"+
		"*%s*?\n\n"+
		"⚠️ Это действие нельзя отменить!", res.DayTypeName)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.YesDelete,
				fmt.Sprintf("workout_delete_%d", res.WorkoutID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.NoCancel,
				fmt.Sprintf("workout_show_progress_%d", res.WorkoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowDeleteWorkout(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "✅ Тренировка успешно удалена!")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.MyWorkouts, "workout_show_my"),
		),
	)
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowMy(chatID int64, res *dto.ShowMyWorkoutsResult) {
	offset, limit, count := res.Pagination.Offset, res.Pagination.Limit, res.Pagination.Total

	var rows [][]tgbotapi.InlineKeyboardButton
	text := fmt.Sprintf("<b>%s</b> (%d-%d из %d):\n\n", messages.MyWorkouts, offset+1, min(offset+limit, count), count)
	for i, workout := range res.Items {

		text += fmt.Sprintf("%d. <u>%s</u> %s\n   %s\n\n",
			i+1+offset, workout.Name, workout.Status, workout.StartedAt)

		// buttons
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}
		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s %d", workout.Name, i+1+offset),
				fmt.Sprintf("workout_show_progress_%d", workout.ID)))
	}

	text += "<b>Выберите тренировку для просмотра:</b>"

	rows = append(rows, []tgbotapi.InlineKeyboardButton{})
	fmt.Println("offset", offset, "limit", limit, "count", count)
	if offset >= limit {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("⬅️ Предыдущие",
			fmt.Sprintf("workout_show_my_%d", offset-limit)))
	}
	if offset+limit < int(count) {
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("➡️ Следующие",
			fmt.Sprintf("workout_show_my_%d", offset+limit)))
	} else {
		rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		rows[len(rows)-1] = append(rows[len(rows)-1], tgbotapi.NewInlineKeyboardButtonData("🔙 В начало", "workout_show_my"))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowStats(chatID int64, res *dto.WorkoutStatistic) {
	dayType := res.DayType
	workoutDay := res.WorkoutDay

	completedExercises := res.CompletedExercises
	totalWeight := res.TotalWeight
	totalTime := res.CardioTime

	exerciseTypesMap := res.ExerciseMap
	exerciseWeightMap := res.ExerciseWeightMap
	exerciseTimeMap := res.ExerciseTimeMap

	var text strings.Builder
	text.WriteString(messages.WorkoutStats + fmt.Sprintf(": %s\n\n", dayType.Name))

	if workoutDay.EndedAt != nil {
		text.WriteString(messages.WorkoutTime + fmt.Sprintf(": %s\n", utils.BetweenTimes(workoutDay.StartedAt, workoutDay.EndedAt)))
	}
	text.WriteString(fmt.Sprintf("<b>%s</b>: %s\n\n", messages.WorkoutDate, workoutDay.StartedAt.Add(3*time.Hour).Format("02.01.2006 15:04")))

	for _, exercise := range workoutDay.Exercises {
		if exercise.CompletedSets() == 0 {
			continue
		}

		exerciseObj, ok := exerciseTypesMap[exercise.ExerciseTypeID]
		if !ok {
			continue
		}

		exerciseWeight, ok := exerciseWeightMap[exercise.ID]
		if !ok {
			continue
		}

		exerciseTime, ok := exerciseTimeMap[exercise.ID]
		if !ok {
			continue
		}

		lastSet := exercise.Sets[len(exercise.Sets)-1]
		text.WriteString(fmt.Sprintf("• <b>%s:</b> \n", exerciseObj.Name))
		if lastSet.GetRealReps() > 0 {
			text.WriteString(fmt.Sprintf("  • Выполнено: %d из %d подходов\n", exercise.CompletedSets(), len(exercise.Sets)))
			text.WriteString(fmt.Sprintf("  • Рабочий вес: %d * %.0f кг \n", lastSet.GetRealReps(), lastSet.GetRealWeight()))
			text.WriteString(fmt.Sprintf("  • Общий вес: %.0f кг \n\n", exerciseWeight))
		} else if lastSet.GetRealMinutes() > 0 {
			text.WriteString(fmt.Sprintf("  • Общее время: %d минут \n\n", exerciseTime))
		}
	}

	text.WriteString(messages.Summary + "\n")
	text.WriteString(fmt.Sprintf("• Упражнений: %d/%d\n", completedExercises, len(workoutDay.Exercises)))
	if totalWeight > 0 {
		text.WriteString(fmt.Sprintf("• Общий тоннаж: %.0f кг\n", totalWeight))
	}
	if totalTime > 0 {
		text.WriteString(fmt.Sprintf("• Общее время: %d минут\n", totalTime))
	}
	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func (p *Presenter) ShowConfirmFinish(chatID, workoutID int64, res *dto.ConfirmFinishWorkout) {
	dayType := res.DayType

	text := fmt.Sprintf("🏁 *Завершение тренировки*\n\n"+
		"Вы уверены, что хотите завершить тренировку:\n"+
		"*%s*?\n\n"+
		"После завершения вы сможете просмотреть статистику, "+
		"но не сможете добавлять новые подходы.", dayType.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, завершить",
				fmt.Sprintf("workout_finish_%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, продолжить",
				fmt.Sprintf("exercise_show_current_session_%d", workoutID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowByUserID(chatID int64, res *dto.ShowWorkoutByUserID) {
	user := res.User
	workouts := res.Workouts

	text := fmt.Sprintf("📋 <b>Тренировки пользователя '%s'</b>\n\n", user.ShortName())
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
		date := workout.StartedAt.Add(3 * time.Hour).Format("02.01.2006 15:04")

		dayType := workout.WorkoutDayType

		text += fmt.Sprintf("%d. <b>%s</b> %s\n   📆️ %s\n\n",
			i+1, dayType.Name, status, date)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	p.bot.Send(msg)
}

func (p *Presenter) ShowCreateWorkoutMenu(chatID int64, program *models.WorkoutProgram) {
	text := "*Выберите день тренировки:*"

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i, day := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(day.Name, fmt.Sprintf("workout_create_%d", day.ID)),
		)
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{})

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = constants.MarkdownParseMode
	p.bot.Send(msg)
}
