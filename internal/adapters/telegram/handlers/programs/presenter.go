package programs

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

func (p *Presenter) ShowProgramManageDialog(chatID int64, result *dto.GetAllPrograms) {
	user := result.User
	programs := result.Programs

	text := &bytes.Buffer{}
	text.WriteString("*Ваши программы:*\n\n")

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, program := range programs {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}

		if program.ID == *user.ActiveProgramID {
			text.WriteString(fmt.Sprintf("• 🟢 *%s* \n  📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		} else {
			text.WriteString(fmt.Sprintf("• *%s* \n 📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		}

		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(program.Name, fmt.Sprintf("program_edit_%d", program.ID)))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новую", "program_create")))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
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

func (p *Presenter) ShowEditDialog(chatID int64, res *dto.GetProgram) {
	program := res.Program
	exerciseTypesMap := res.ExerciseTypesMap

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	text := &bytes.Buffer{}

	text.WriteString(fmt.Sprintf("<b>Программа: %s</b>\n\n", program.Name))
	text.WriteString("<b>Список дней:</b>\n\n")
	for i, dayType := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(dayType.Name, fmt.Sprintf("day_type_edit_%d", dayType.ID)),
		)

		text.WriteString(fmt.Sprintf("<b>%d. %s</b>\n", i+1, dayType.Name))
		text.WriteString(fmt.Sprintf("%s \n\n", formatPreset(dayType.Preset, exerciseTypesMap)))
	}
	text.WriteString("<b>Выберите день, в который хотите добавить упражнения:</b>")

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить день", fmt.Sprintf("change_day_name_%d", program.ID)),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Переименовать", fmt.Sprintf("change_name_of_program_%d", program.ID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👑 Выбрать текущей", fmt.Sprintf("program_change_%d", program.ID)),
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить", fmt.Sprintf("program_confirm_delete_%d", program.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ConfirmDeleteDialog(chatID int64, res *dto.GetProgram) {
	program := res.Program
	text := fmt.Sprintf("🗑️ *Удаление программы*\n\n"+
		"Вы уверены, что хотите удалить программу:\n"+
		"*%s*?\n\n"+
		"⚠️ Это действие нельзя отменить!", program.Name)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить",
				fmt.Sprintf("program_delete_%d", program.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Нет, отмена",
				fmt.Sprintf("program_edit_%d", program.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func formatPreset(preset string, exerciseTypesMap map[int64]models.ExerciseType) string {
	exercises := utils.SplitPreset(preset)
	buffer := &bytes.Buffer{}
	for _, ex := range exercises {
		exerciseType, ok := exerciseTypesMap[ex.ID]
		if !ok {
			continue
		}
		buffer.WriteString(fmt.Sprintf("• <u>%s</u>\n", exerciseType.Name))
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
		buffer.WriteString("\n")
	}
	return buffer.String()
}
