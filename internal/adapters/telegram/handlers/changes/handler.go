package changes

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/changes/userstatemachine"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises/presenter"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/groups"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/session"
	"github.com/SaenkoDmitry/training-tg-bot/internal/constants"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	daytypeusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/daytypes"
	exercisecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exercises"
	programusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/programs"
	setusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/sets"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type Handler struct {
	presenter         *Presenter
	commonPresenter   *common.Presenter
	exercisePresenter *presenter.Presenter
	programPresenter  *programs.Presenter

	userStatesMachine *userstatemachine.UserStatesMachine

	showCurrentSessionUC *session.ShowCurrentExerciseSessionUseCase
	changeNextSetUC      *setusecases.UpdateNextUseCase
	manageProgramUC      *programusecases.FindAllByUserUseCase
	renameProgramUC      *programusecases.RenameUseCase
	getAllGroupsUC       *groups.GetAllUseCase
	dayTypesCreateUC     *daytypeusecases.CreateUseCase
	exerciseTypeListUC   *exercisecases.ExerciseTypeListUseCase
	getProgramUC         *programusecases.GetUseCase
	dayTypesUpdateUC     *daytypeusecases.UpdateUseCase
	dayTypeGetUC         *daytypeusecases.GetUseCase
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	showCurrentSessionUC *session.ShowCurrentExerciseSessionUseCase,
	changeNextSetUC *setusecases.UpdateNextUseCase,
	manageProgramUC *programusecases.FindAllByUserUseCase,
	renameProgramUC *programusecases.RenameUseCase,
	getAllGroupsUC *groups.GetAllUseCase,
	dayTypesCreateUC *daytypeusecases.CreateUseCase,
	dayTypesUpdateUC *daytypeusecases.UpdateUseCase,
	dayTypeGetUC *daytypeusecases.GetUseCase,
	exerciseTypeListUC *exercisecases.ExerciseTypeListUseCase,
	editProgramUC *programusecases.GetUseCase,
) *Handler {
	return &Handler{
		presenter:            NewPresenter(bot),
		commonPresenter:      common.NewPresenter(bot),
		exercisePresenter:    presenter.NewPresenter(bot),
		programPresenter:     programs.NewPresenter(bot),
		userStatesMachine:    userstatemachine.New(),
		showCurrentSessionUC: showCurrentSessionUC,
		changeNextSetUC:      changeNextSetUC,
		manageProgramUC:      manageProgramUC,
		renameProgramUC:      renameProgramUC,
		getAllGroupsUC:       getAllGroupsUC,
		dayTypesCreateUC:     dayTypesCreateUC,
		dayTypesUpdateUC:     dayTypesUpdateUC,
		dayTypeGetUC:         dayTypeGetUC,
		exerciseTypeListUC:   exerciseTypeListUC,
		getProgramUC:         editProgramUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "change_reps_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_reps_ex_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_reps_%d", exerciseID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterNewReps)

	case strings.HasPrefix(data, "change_weight_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_weight_ex_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_weight_%d", exerciseID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterNewWeight)

	case strings.HasPrefix(data, "change_minutes_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_minutes_ex_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_minutes_%d", exerciseID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterNewTime)

	case strings.HasPrefix(data, "change_meters_ex_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_meters_ex_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_meters_%d", exerciseID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterNewMeters)

	case strings.HasPrefix(data, "change_day_name_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_day_name_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_name_for_program_%d", programID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterWorkoutDayName)

	case strings.HasPrefix(data, "change_name_of_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "change_name_of_program_"), 10, 64)
		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_program_name_%d", programID))
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.EnterNewProgramName)

	case strings.HasPrefix(data, "change_program_day_add_exercise_"):
		parts := strings.Split(strings.TrimPrefix(data, "change_program_day_add_exercise_"), "_")
		if len(parts) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(parts[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(parts[1], 10, 64)

		h.userStatesMachine.SetValue(chatID, fmt.Sprintf("awaiting_day_preset_%d_%d", dayTypeID, exerciseTypeID))

		text := messages.EnterPreset
		if exerciseTypesResult, err := h.exerciseTypeListUC.Execute(); err == nil {
			exerciseTypeUnits := constants.RepsUnit + "," + constants.WeightUnit
			for _, ex := range exerciseTypesResult.ExerciseTypes {
				if ex.ID == exerciseTypeID && ex.Units != "" {
					exerciseTypeUnits = ex.Units
					break
				}
			}
			text += fmt.Sprintf("\n\n<b>Подсказка:</b> для вашего упражнения следует выбрать <b>%s</b> !", exerciseTypeUnits)
		}
		h.commonPresenter.SendSimpleHtmlMessage(chatID, text)
	}
}

func (h *Handler) RouteMessage(chatID int64, text string) {
	state, exists := h.userStatesMachine.GetValue(chatID)
	if !exists {
		return
	}

	switch {
	case strings.HasPrefix(state, "awaiting_reps_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_reps_"), 10, 64)
		newReps, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.IncorrectFormatReps)
			return
		}
		workoutID := h.updateNextSet(chatID, exerciseID, &dto.NewSet{NewReps: newReps})
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.RepsUpdated)
		if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(workoutID); sessionErr == nil {
			h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
		}

	case strings.HasPrefix(state, "awaiting_weight_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_weight_"), 10, 64)
		newWeight, err := strconv.ParseFloat(text, 32)
		if err != nil {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.IncorrectFormatWeight)
			return
		}
		workoutID := h.updateNextSet(chatID, exerciseID, &dto.NewSet{NewWeight: newWeight})
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.WeightUpdated)
		if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(workoutID); sessionErr == nil {
			h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
		}

	case strings.HasPrefix(state, "awaiting_minutes_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_minutes_"), 10, 64)
		newMinutes, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.IncorrectFormatMinutes)
			return
		}
		workoutID := h.updateNextSet(chatID, exerciseID, &dto.NewSet{NewMinutes: newMinutes})
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.MinutesUpdated)
		if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(workoutID); sessionErr == nil {
			h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
		}

	case strings.HasPrefix(state, "awaiting_meters_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_meters_"), 10, 64)
		newMeters, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.IncorrectFormatMeters)
			return
		}
		workoutID := h.updateNextSet(chatID, exerciseID, &dto.NewSet{NewMeters: newMeters})
		h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.MetersUpdated)
		if sessionResult, sessionErr := h.showCurrentSessionUC.Execute(workoutID); sessionErr == nil {
			h.exercisePresenter.ShowCurrentSession(chatID, sessionResult)
		}

	case strings.HasPrefix(state, "awaiting_program_name_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_program_name_"), 10, 64)
		err := h.renameProgramUC.Execute(programID, text)
		if err != nil {
			h.commonPresenter.HandleInternalError(err, chatID, h.changeNextSetUC.Name())
			return
		}
		if res, manageErr := h.manageProgramUC.Execute(chatID); manageErr == nil {
			h.programPresenter.ShowProgramManageDialog(chatID, res)
		}

	case strings.HasPrefix(state, "awaiting_day_name_for_program_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(state, "awaiting_day_name_for_program_"), 10, 64)
		dayTypeID, err := h.dayTypesCreateUC.Execute(programID, text)
		if err != nil {
			h.commonPresenter.HandleInternalError(err, chatID, h.dayTypesCreateUC.Name())
			return
		}
		if res, addErr := h.getAllGroupsUC.Execute(); addErr == nil {
			h.programPresenter.ShowSelectDayTypeDialog(chatID, dayTypeID, res)
		}

	case strings.HasPrefix(state, "awaiting_day_preset_"):
		text = strings.ToLower(text)

		exerciseTypeListResult, err := h.exerciseTypeListUC.Execute()
		if err != nil {
			h.commonPresenter.HandleInternalError(err, chatID, h.exerciseTypeListUC.Name())
			return
		}

		// parse dayTypeID and exerciseTypeID
		parts := strings.Split(strings.TrimPrefix(state, "awaiting_day_preset_"), "_")
		if len(parts) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(parts[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(parts[1], 10, 64)

		found := false
		var exerciseType models.ExerciseType
		for _, exType := range exerciseTypeListResult.ExerciseTypes {
			if exType.ID == exerciseTypeID {
				exerciseType = exType
				found = true
				break
			}
		}
		if !found {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, "Не найдено упражнение")
			return
		}

		textArr := strings.Split(text, ":")
		if len(textArr) != 2 {
			h.sendIncorrectPresetMsg(chatID, exerciseType.Units)
			return
		}

		preset := textArr[1]

		units, valid := utils.SplitUnits(textArr[0])
		if !valid {
			h.sendIncorrectPresetMsg(chatID, exerciseType.Units)
			return
		}
		exUnits, _ := utils.SplitUnits(exerciseType.Units)

		if !utils.EqualArrays(exUnits, units) {
			h.sendIncorrectPresetMsg(chatID, exerciseType.Units)
			return
		}
		presetSetLen := 1
		if strings.Contains(preset, "*") {
			presetSetLen = 2
		}
		if len(exUnits) != presetSetLen {
			h.sendIncorrectPresetMsg(chatID, exerciseType.Units)
			return
		}

		if !utils.IsValidPreset(preset) {
			h.sendIncorrectPresetMsg(chatID, exerciseType.Units)
			return
		}

		dayType, err := h.dayTypeGetUC.Execute(dayTypeID)
		if err != nil {
			h.commonPresenter.HandleInternalError(err, chatID, h.dayTypeGetUC.Name())
			return
		}
		if dayType.Preset != "" {
			dayType.Preset += ";"
		}

		dayType.Preset += fmt.Sprintf("%d:[%s]", exerciseTypeID, preset)

		if updateErr := h.dayTypesUpdateUC.Execute(dayType); updateErr != nil {
			h.commonPresenter.HandleInternalError(err, chatID, h.dayTypesUpdateUC.Name())
			return
		}

		if editResult, editErr := h.getProgramUC.Execute(dayType.WorkoutProgramID); editErr == nil {
			h.programPresenter.ViewProgram(chatID, editResult)
		}
	}
}

func (h *Handler) sendIncorrectPresetMsg(chatID int64, expectedUnits string) {
	h.commonPresenter.SendSimpleHtmlMessage(chatID, "❌ Неверный формат !\n\n"+messages.EnterPreset+
		fmt.Sprintf("\n\n<b>Подсказка:</b> для вашего упражнения следует выбрать <b>%s</b> !", expectedUnits))
}

func (h *Handler) updateNextSet(chatID, exerciseID int64, newSet *dto.NewSet) int64 {
	workoutID, err := h.changeNextSetUC.Execute(exerciseID, newSet)
	if err != nil {
		if errors.Is(err, session.NotFoundExerciseErr) {
			return 0
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.changeNextSetUC.Name())
		return 0
	}
	h.userStatesMachine.Clear(chatID)
	return workoutID
}
