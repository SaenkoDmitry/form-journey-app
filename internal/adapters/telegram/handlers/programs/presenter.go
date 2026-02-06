package programs

import (
	"bytes"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
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
	text.WriteString(fmt.Sprintf("<b>%s:</b>\n\n", messages.MyPrograms))

	var rows [][]tgbotapi.InlineKeyboardButton
	for i, program := range programs {
		if i%2 == 0 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{})
		}

		if program.ID == *user.ActiveProgramID {
			text.WriteString(fmt.Sprintf("•  👉 <b>%s</b> \n  📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		} else {
			text.WriteString(fmt.Sprintf("•    <b>%s</b> \n 📅 %s\n\n", program.Name, program.CreatedAt.Format("02.01.2006 15:04")))
		}

		rows[len(rows)-1] = append(rows[len(rows)-1],
			tgbotapi.NewInlineKeyboardButtonData(program.Name, fmt.Sprintf("program_view_%d", program.ID)))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Добавить новую", "program_create"),
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, "/menu"),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
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

func (p *Presenter) ViewAllDays(chatID int64, res *dto.GetProgram) {
	program := res.Program

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	text := &bytes.Buffer{}

	text.WriteString(fmt.Sprintf("<b>Программа:</b> %s\n", program.Name))

	for i, dayType := range program.DayTypes {
		if i%2 == 0 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow())
		}
		buttons[len(buttons)-1] = append(buttons[len(buttons)-1],
			tgbotapi.NewInlineKeyboardButtonData(dayType.Name, fmt.Sprintf("day_type_view_%d", dayType.ID)),
		)
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, fmt.Sprintf("program_view_%d", program.ID)),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	text.WriteString("\n<b>Выберите день для просмотра:</b>")

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.ParseMode = constants.HtmlParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}

func (p *Presenter) ViewProgram(chatID int64, res *dto.GetProgram) {
	program := res.Program

	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	text := &bytes.Buffer{}

	text.WriteString(fmt.Sprintf("<b>Программа: %s</b>\n\n", program.Name))
	text.WriteString("<b>Список дней:</b>\n\n")

	for i, dayType := range program.DayTypes {
		text.WriteString(fmt.Sprintf("<b>%d.</b> %s\n", i+1, dayType.Name))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👑 Выбрать текущей", fmt.Sprintf("program_change_%d", program.ID)),
		tgbotapi.NewInlineKeyboardButtonData("🎟️ Переименовать", fmt.Sprintf("change_name_of_program_%d", program.ID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("➕ Еще день", fmt.Sprintf("change_day_name_%d", program.ID)),
		tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить всю", fmt.Sprintf("program_confirm_delete_%d", program.ID)),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🕹️ Управление днями", fmt.Sprintf("program_view_all_days_%d", program.ID)),
		tgbotapi.NewInlineKeyboardButtonData(messages.BackTo, "program_management"),
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
			tgbotapi.NewInlineKeyboardButtonData(messages.YesDelete,
				fmt.Sprintf("program_delete_%d", program.ID)),
			tgbotapi.NewInlineKeyboardButtonData(messages.NoCancel,
				fmt.Sprintf("program_view_%d", program.ID)),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = constants.MarkdownParseMode
	msg.ReplyMarkup = keyboard
	p.bot.Send(msg)
}
