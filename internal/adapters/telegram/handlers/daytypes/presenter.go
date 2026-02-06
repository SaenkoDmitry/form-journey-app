package daytypes

import (
	"bytes"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Presenter struct {
	bot *tgbotapi.BotAPI
}

func NewPresenter(bot *tgbotapi.BotAPI) *Presenter {
	return &Presenter{bot: bot}
}

func (p *Presenter) ShowSelectDayTypeDialog(chatID int64, dayTypeID int64, res *dto.ExerciseGroupTypeList) {
	groups := res.Groups

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for i, group := range groups {
		if i%3 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(group.Name, fmt.Sprintf("exercise_select_for_program_day_%d_%s", dayTypeID, group.Code)),
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, messages.SelectGroupOfMuscle)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ShowConfirmDelete(chatID int64, res *models.WorkoutDayType) {
	text := fmt.Sprintf("🗑️ <b>Удаление тренировочного дня из программы</b>\n\n"+
		"Вы уверены, что хотите удалить день:\n"+
		"<b>%s</b>?\n\n"+
		"⚠️ Это действие нельзя отменить!", res.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(messages.YesDelete,
				fmt.Sprintf("day_type_delete_%d", res.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.NoCancel,
				fmt.Sprintf("program_view_%d", res.WorkoutProgramID)),
		),
	)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ViewDayType(chatID int64, res *models.WorkoutDayType, programsResult *dto.GetProgram) {
	program := programsResult.Program
	exerciseTypesMap := programsResult.ExerciseTypesMap
	daytypeID := res.ID

	text := &bytes.Buffer{}
	text.WriteString(fmt.Sprintf("<b>День:</b> %s\n\n", res.Name))
	text.WriteString(fmt.Sprintf("%s \n\n", formatPreset(res.Preset, exerciseTypesMap)))

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(

		tgbotapi.NewInlineKeyboardButtonData("✏️️ Добавить упражнение", fmt.Sprintf("day_type_edit_%d", daytypeID)),
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", fmt.Sprintf("day_type_confirm_delete_%d", daytypeID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, fmt.Sprintf("program_view_all_days_%d", program.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func formatPreset(preset string, exerciseTypesMap map[int64]models.ExerciseType) string {
	exercises := utils.SplitPreset(preset)
	buffer := &bytes.Buffer{}
	for i, ex := range exercises {
		exerciseType, ok := exerciseTypesMap[ex.ID]
		if !ok {
			continue
		}
		buffer.WriteString(fmt.Sprintf("• <b>%d.</b> <u>%s</u>\n", i+1, exerciseType.Name))
		buffer.WriteString(fmt.Sprintf("    • "))
		for i, set := range ex.Sets {
			if i > 0 {
				buffer.WriteString(", ")
			}
			if set.Minutes > 0 {
				buffer.WriteString(fmt.Sprintf("%d мин", set.Minutes))
			} else {
				buffer.WriteString(fmt.Sprintf("%d * %.0f кг", set.Reps, set.Weight))
			}
		}
		buffer.WriteString("\n\n")
	}
	return buffer.String()
}
