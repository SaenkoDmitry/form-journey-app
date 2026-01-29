package exercises

import (
	"errors"
	"fmt"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	presenter2 "github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises/presenter"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/workouts"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/exercises"
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/groups"
	sessionusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type Handler struct {
	presenter       *presenter2.Presenter
	commonPresenter *common.Presenter

	getGroupUC     *groups.GetUseCase
	getAllGroupsUC *groups.GetAllUseCase

	moveToExerciseSessionUC *sessionusecases.MoveToUseCase

	showCurrentSessionUC    *sessionusecases.ShowCurrentExerciseSessionUseCase
	findTypesByGroupUC      *exercises.FindTypesByGroupUseCase
	confirmDeleteExerciseUC *exercises.ConfirmDeleteUseCase
	deleteExerciseUC        *exercises.DeleteUseCase
	getExerciseUC           *exercises.GetUseCase
	createExerciseUC        *exercises.CreateUseCase
	workoutsHandler         *workouts.Handler
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	showCurrentExerciseUC *sessionusecases.ShowCurrentExerciseSessionUseCase,
	getGroupUC *groups.GetUseCase,
	findTypesByGroupUC *exercises.FindTypesByGroupUseCase,
	confirmDeleteExerciseUC *exercises.ConfirmDeleteUseCase,
	deleteExerciseUC *exercises.DeleteUseCase,
	moveToExerciseSessionUC *sessionusecases.MoveToUseCase,
	getExerciseUC *exercises.GetUseCase,
	showGroupTypeListUC *groups.GetAllUseCase,
	createExerciseUC *exercises.CreateUseCase,
	workoutsHandler *workouts.Handler,
) *Handler {
	return &Handler{
		presenter:               presenter2.NewPresenter(bot),
		commonPresenter:         common.NewPresenter(bot),
		showCurrentSessionUC:    showCurrentExerciseUC,
		getGroupUC:              getGroupUC,
		findTypesByGroupUC:      findTypesByGroupUC,
		confirmDeleteExerciseUC: confirmDeleteExerciseUC,
		deleteExerciseUC:        deleteExerciseUC,
		moveToExerciseSessionUC: moveToExerciseSessionUC,
		getExerciseUC:           getExerciseUC,
		getAllGroupsUC:          showGroupTypeListUC,
		createExerciseUC:        createExerciseUC,
		workoutsHandler:         workoutsHandler,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "exercise_show_current_session_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_show_current_session_"), 10, 64)
		h.ShowCurrentExerciseSession(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_move_to_prev_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_move_to_prev_"), 10, 64)
		h.MoveToPrevExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_move_to_next_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_move_to_next_"), 10, 64)
		h.MoveToNextExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_show_hint_"):
		parts := strings.Split(strings.TrimPrefix(data, "exercise_show_hint_"), "_")
		workoutID, _ := strconv.ParseInt(parts[0], 10, 64)
		exerciseTypeID, _ := strconv.ParseInt(parts[1], 10, 64)
		h.showExerciseHint(chatID, workoutID, exerciseTypeID)

	case strings.HasPrefix(data, "exercise_show_info_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_show_info_"), 10, 64)
		h.showExerciseHint(chatID, 0, exerciseID)

	case strings.HasPrefix(data, "exercise_add_for_current_workout_"):
		workoutDayID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_add_for_current_workout_"), 10, 64)
		h.addExercise(chatID, workoutDayID)

	case strings.HasPrefix(data, "exercise_select_for_current_workout_"):
		text := strings.TrimPrefix(data, "exercise_select_for_current_workout_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutDayID, _ := strconv.ParseInt(arr[0], 10, 64)
			code := arr[1]
			h.selectExerciseForCurrentWorkout(chatID, workoutDayID, code)
		}

	case strings.HasPrefix(data, "exercise_select_for_program_day_"):
		parts := strings.Split(strings.TrimPrefix(data, "exercise_select_for_program_day_"), "_")
		if len(parts) < 2 {
			return
		}
		dayTypeID, _ := strconv.ParseInt(parts[0], 10, 64)
		exerciseGroupCode := parts[1]
		h.selectExerciseForProgramDay(chatID, dayTypeID, exerciseGroupCode)

	case strings.HasPrefix(data, "exercise_add_specific_for_current_workout_"):
		text := strings.TrimPrefix(data, "exercise_add_specific_for_current_workout_")
		if arr := strings.Split(text, "_"); len(arr) == 2 {
			workoutID, _ := strconv.ParseInt(arr[0], 10, 64)
			internalExerciseID, _ := strconv.ParseInt(arr[1], 10, 64)
			h.addSpecificExerciseForCurrentWorkout(chatID, workoutID, internalExerciseID)
		}

	case strings.HasPrefix(data, "exercise_confirm_delete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_confirm_delete_"), 10, 64)
		h.confirmDeleteExercise(chatID, exerciseID)

	case strings.HasPrefix(data, "exercise_delete_"):
		exerciseID, _ := strconv.ParseInt(strings.TrimPrefix(data, "exercise_delete_"), 10, 64)
		h.deleteExercise(chatID, exerciseID)

	case strings.HasPrefix(data, "exercise_show_all_groups"):
		h.showAllGroups(chatID)

	case strings.HasPrefix(data, "exercise_show_list_"):
		groupCode := strings.TrimPrefix(data, "exercise_show_list_")
		h.showAllExercisesByGroup(chatID, groupCode)
	}
}

func (h *Handler) ShowCurrentExerciseSession(chatID, workoutID int64) {
	res, err := h.showCurrentSessionUC.Execute(workoutID)
	if err != nil {
		if errors.Is(err, sessionusecases.NoExercisesErr) {
			h.presenter.ShowNoExercises(chatID)
			return
		}
		if errors.Is(err, sessionusecases.NotFoundExerciseErr) {
			h.presenter.ShowNotFoundExercise(chatID)
			return
		}
		return
	}
	h.presenter.ShowCurrentSession(chatID, res)
}

func (h *Handler) selectExerciseForProgramDay(chatID int64, dayTypeID int64, exerciseGroupCode string) {
	group, err := h.getGroupUC.Execute(exerciseGroupCode)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getGroupUC.Name())
		return
	}

	result, err := h.findTypesByGroupUC.Execute(exerciseGroupCode)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findTypesByGroupUC.Name())
		return
	}

	h.presenter.ShowSelectExerciseForProgramDayDialog(chatID, dayTypeID, group, result.ExerciseTypes)
}

func (h *Handler) confirmDeleteExercise(chatID int64, exerciseID int64) {
	res, err := h.confirmDeleteExerciseUC.Execute(exerciseID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.confirmDeleteExerciseUC.Name())
		return
	}

	h.presenter.ShowConfirmDeleteDialog(chatID, res)
}

func (h *Handler) deleteExercise(chatID int64, exerciseID int64) {
	workoutID, err := h.deleteExerciseUC.Execute(exerciseID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.deleteExerciseUC.Name())
		return
	}

	h.ShowCurrentExerciseSession(chatID, workoutID)
}

func (h *Handler) MoveToPrevExercise(chatID int64, workoutID int64) {
	h.moveToExercise(chatID, workoutID, false)
}

func (h *Handler) MoveToNextExercise(chatID int64, workoutID int64) {
	h.moveToExercise(chatID, workoutID, true)
}

func (h *Handler) moveToExercise(chatID int64, workoutID int64, next bool) {
	err := h.moveToExerciseSessionUC.Execute(workoutID, next)
	if err != nil {
		if errors.Is(err, sessionusecases.NoExercisesInWorkout) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, "❌ В тренировке нет упражнений")
			return
		}
		if errors.Is(err, sessionusecases.NoEarlierExercisesInWorkout) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, "❌ Более ранних упражнений в этой тренировке нет")
			h.ShowCurrentExerciseSession(chatID, workoutID)
			return
		}

		if errors.Is(err, sessionusecases.YouCompletedAllExercises) {
			h.presenter.CompleteAllExercises(chatID, workoutID)
			return
		}

		h.commonPresenter.HandleInternalError(err, chatID, h.moveToExerciseSessionUC.Name())
		return
	}
	h.ShowCurrentExerciseSession(chatID, workoutID)
}

func (h *Handler) showExerciseHint(chatID int64, workoutID, exerciseID int64) {
	res, err := h.getExerciseUC.Execute(exerciseID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getExerciseUC.Name())
		return
	}
	h.presenter.ShowHint(chatID, res, workoutID)
}

func (h *Handler) addExercise(chatID int64, workoutID int64) {
	groupsResult, err := h.getAllGroupsUC.Execute()
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getAllGroupsUC.Name())
		return
	}

	h.presenter.AddExerciseDialog(chatID, workoutID, groupsResult.Groups)
}

func (h *Handler) selectExerciseForCurrentWorkout(chatID int64, workoutID int64, code string) {
	group, err := h.getGroupUC.Execute(code)
	if err != nil {
		return
	}

	exerciseTypeResult, err := h.findTypesByGroupUC.Execute(code)
	if err != nil {
		return
	}

	h.presenter.ShowSelectExerciseForCurrentWorkoutDialog(chatID, workoutID, group, exerciseTypeResult.ExerciseTypes)
}

func (h *Handler) addSpecificExerciseForCurrentWorkout(chatID int64, workoutID int64, internalExerciseID int64) {
	res, err := h.createExerciseUC.Execute(workoutID, internalExerciseID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.createExerciseUC.Name())
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, fmt.Sprintf("Упражнение <b>'%s'</b> добавлено! ✅", res.ExerciseObj.Name))

	h.workoutsHandler.ShowProgress(chatID, workoutID)
}

func (h *Handler) showAllGroups(chatID int64) {
	groupsResult, err := h.getAllGroupsUC.Execute()
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getAllGroupsUC.Name())
		return
	}
	h.presenter.ShowAllGroups(chatID, groupsResult.Groups)
}

func (h *Handler) showAllExercisesByGroup(chatID int64, groupCode string) {
	exercisesResult, err := h.findTypesByGroupUC.Execute(groupCode)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findTypesByGroupUC.Name())
		return
	}
	res, err := h.getGroupUC.Execute(groupCode)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getGroupUC.Name())
		return
	}
	h.presenter.ShowAllExercises(chatID, exercisesResult.ExerciseTypes, res.Name)
}
