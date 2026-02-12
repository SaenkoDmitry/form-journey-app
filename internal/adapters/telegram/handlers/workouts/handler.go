package workouts

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	exercisepresenter "github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/handlers/exercises/presenter"
	programusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/programs"
	exerciseusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/session"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	workoutusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/workouts"
)

type Handler struct {
	deleteUC        *workoutusecases.DeleteUseCase
	confirmDeleteUC *workoutusecases.ConfirmDeleteUseCase
	createUC        *workoutusecases.CreateUseCase
	startUC         *workoutusecases.StartUseCase
	showMyUC        *workoutusecases.FindMyUseCase
	showProgressUC  *workoutusecases.ShowProgressUseCase
	finishUC        *workoutusecases.FinishUseCase
	confirmFinishUC *workoutusecases.ConfirmFinishUseCase
	showByUserIDUC  *workoutusecases.FindByUserIDUseCase
	statsUC         *workoutusecases.StatsUseCase

	getByUserProgramUC *programusecases.GetByUserUseCase

	showCurrentExerciseSessionUC *exerciseusecases.ShowCurrentExerciseSessionUseCase

	presenter          *Presenter
	exercisesPresenter *exercisepresenter.Presenter
	commonPresenter    *common.Presenter
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	deleteUC *workoutusecases.DeleteUseCase,
	confirmDeleteUC *workoutusecases.ConfirmDeleteUseCase,
	createUC *workoutusecases.CreateUseCase,
	startUC *workoutusecases.StartUseCase,
	showMyUC *workoutusecases.FindMyUseCase,
	showProgressUC *workoutusecases.ShowProgressUseCase,
	finishUC *workoutusecases.FinishUseCase,
	confirmFinishUC *workoutusecases.ConfirmFinishUseCase,
	ShowCurrentExerciseSessionUC *exerciseusecases.ShowCurrentExerciseSessionUseCase,
	showByUserIDUC *workoutusecases.FindByUserIDUseCase,
	statsUC *workoutusecases.StatsUseCase,
	getByUserProgramUC *programusecases.GetByUserUseCase,
) *Handler {
	return &Handler{
		deleteUC:                     deleteUC,
		confirmDeleteUC:              confirmDeleteUC,
		createUC:                     createUC,
		startUC:                      startUC,
		showMyUC:                     showMyUC,
		showProgressUC:               showProgressUC,
		finishUC:                     finishUC,
		confirmFinishUC:              confirmFinishUC,
		showByUserIDUC:               showByUserIDUC,
		showCurrentExerciseSessionUC: ShowCurrentExerciseSessionUC,
		getByUserProgramUC:           getByUserProgramUC,
		statsUC:                      statsUC,

		presenter:          NewPresenter(bot),
		commonPresenter:    common.NewPresenter(bot),
		exercisesPresenter: exercisepresenter.NewPresenter(bot),
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case data == "workout_show_my":
		h.showMy(chatID, 0)

	case strings.HasPrefix(data, "workout_show_my_"):
		offset, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_my_"), 10, 64)
		h.showMy(chatID, int(offset))

	case strings.HasPrefix(data, "workout_create_"):
		dayTypeID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_create_"), 10, 64)
		if workoutID := h.create(chatID, dayTypeID); workoutID != 0 {
			h.ShowProgress(chatID, workoutID, true)
		}

	case strings.HasPrefix(data, "workout_start_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_start_"), 10, 64)
		h.start(chatID, workoutID)

	case strings.HasPrefix(data, "workout_show_progress_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_progress_"), 10, 64)
		h.ShowProgress(chatID, workoutID, true)

	case strings.HasPrefix(data, "workout_confirm_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_confirm_delete_"), 10, 64)
		h.confirmDelete(chatID, workoutID)

	case strings.HasPrefix(data, "workout_delete_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_delete_"), 10, 64)
		h.delete(chatID, workoutID)

	case strings.HasPrefix(data, "workout_confirm_finish_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_confirm_finish_"), 10, 64)
		h.confirmFinish(chatID, workoutID)

	case strings.HasPrefix(data, "workout_finish_"):
		workoutID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_finish_"), 10, 64)
		h.finish(chatID, workoutID)

	case strings.HasPrefix(data, "workout_show_by_user_id_"):
		userID, _ := strconv.ParseInt(strings.TrimPrefix(data, "workout_show_by_user_id_"), 10, 64)
		h.showByUserID(chatID, userID)
	}
}

func (h *Handler) RouteMessage(chatID int64, text string) {
	switch {
	case strings.EqualFold(text, "/workouts/start"):
		h.showWorkoutTypeMenu(chatID)

	case strings.EqualFold(text, "/workouts"):
		h.showMy(chatID, 0)
	}
}

func (h *Handler) showWorkoutTypeMenu(chatID int64) {
	program, err := h.getByUserProgramUC.Execute(chatID)
	if err != nil {
		return
	}
	if len(program.DayTypes) == 0 {
		h.commonPresenter.SendSimpleHtmlMessage(chatID, "Добавьте тренировочные дни в программу через '⚙️ Настройки'")
		return
	}

	h.presenter.ShowCreateWorkoutMenu(chatID, program)
}

func (h *Handler) confirmDelete(chatID int64, workoutID int64) {
	res, err := h.confirmDeleteUC.Execute(workoutID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.confirmDeleteUC.Name())
		return
	}
	h.presenter.ShowConfirmDeleteWorkout(chatID, res)
}

func (h *Handler) delete(chatID int64, workoutID int64) {
	err := h.deleteUC.Execute(workoutID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.deleteUC.Name())
		return
	}
	h.presenter.ShowDeleteWorkout(chatID)
}

func (h *Handler) create(chatID int64, dayTypeID int64) int64 {
	res, err := h.createUC.Execute(chatID, dayTypeID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.createUC.Name())
		return 0
	}
	h.presenter.WorkoutCreated(chatID)
	return res.WorkoutID
}

func (h *Handler) ShowProgress(chatID int64, workoutID int64, needShowButtons bool) {
	workoutProgress, err := h.showProgressUC.Execute(workoutID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.showProgressUC.Name())
		return
	}
	statsResult, err := h.statsUC.Execute(workoutID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.statsUC.Name())
		return
	}
	h.presenter.ShowWorkoutProgress(chatID, workoutProgress, statsResult, needShowButtons)
}

func (h *Handler) start(chatID int64, workoutID int64) {
	if _, err := h.startUC.Execute(workoutID); err != nil {
		switch {
		case errors.Is(err, workoutusecases.NotFoundSpecificErr):
			h.presenter.ShowNotFoundSpecific(chatID)
		case errors.Is(err, workoutusecases.AlreadyCompletedErr):
			h.presenter.ShowAlreadyCompleted(chatID)
		default:
			h.commonPresenter.HandleInternalError(err, chatID, h.createUC.Name())
		}
		return
	}
	if res, showErr := h.showCurrentExerciseSessionUC.Execute(workoutID); showErr == nil {
		h.exercisesPresenter.ShowCurrentSession(chatID, res)
	}
}

const showWorkoutsLimit = 4

func (h *Handler) showMy(chatID int64, offset int) {
	res, err := h.showMyUC.Execute(chatID, offset, showWorkoutsLimit)
	if err != nil {
		if errors.Is(err, workoutusecases.NotFoundAllErr) {
			h.presenter.ShowNotFoundAll(chatID)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.showMyUC.Name())
		return
	}
	h.presenter.ShowMy(chatID, res)
}

func (h *Handler) confirmFinish(chatID int64, workoutID int64) {
	res, err := h.confirmFinishUC.Execute(workoutID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.showMyUC.Name())
		return
	}
	h.presenter.ShowConfirmFinish(chatID, workoutID, res)
}

func (h *Handler) finish(chatID int64, workoutID int64) {
	if _, err := h.finishUC.Execute(workoutID); err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.showMyUC.Name())
		return
	}
	h.ShowProgress(chatID, workoutID, false)
}

func (h *Handler) showByUserID(chatID int64, userID int64) {
	res, err := h.showByUserIDUC.Execute(userID)
	if err != nil {
		if errors.Is(err, workoutusecases.EmptyWorkoutsErr) {
			h.presenter.ShowNotFoundAllForUser(chatID, res.User)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.showMyUC.Name())
		return
	}
	h.presenter.ShowByUserID(chatID, res)
}
