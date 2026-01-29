package programs

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	"github.com/SaenkoDmitry/training-tg-bot/internal/messages"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	programusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/programs"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	deleteProgramUC         *programusecases.DeleteUseCase
	createProgramUC         *programusecases.CreateUseCase
	activateProgramUC       *programusecases.ActivateUseCase
	getProgramUC            *programusecases.GetUseCase
	findAllProgramsByUserUC *programusecases.FindAllByUserUseCase
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	deleteProgramUC *programusecases.DeleteUseCase,
	createProgramUC *programusecases.CreateUseCase,
	activateProgramUC *programusecases.ActivateUseCase,
	editProgramUC *programusecases.GetUseCase,
	manageProgramUC *programusecases.FindAllByUserUseCase,
) *Handler {
	return &Handler{
		presenter:               NewPresenter(bot),
		commonPresenter:         common.NewPresenter(bot),
		deleteProgramUC:         deleteProgramUC,
		createProgramUC:         createProgramUC,
		activateProgramUC:       activateProgramUC,
		getProgramUC:            editProgramUC,
		findAllProgramsByUserUC: manageProgramUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {
	case strings.HasPrefix(data, "program_create"):
		h.createProgram(chatID)

	case strings.HasPrefix(data, "program_management"):
		h.programManagement(chatID)

	case strings.HasPrefix(data, "program_edit_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_edit_"), 10, 64)
		h.editProgram(chatID, programID)

	case strings.HasPrefix(data, "program_change_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_change_"), 10, 64)
		h.activateProgram(chatID, programID)

	case strings.HasPrefix(data, "program_confirm_delete_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_confirm_delete_"), 10, 64)
		h.confirmDeleteProgram(chatID, programID)

	case strings.HasPrefix(data, "program_delete_"):
		programID, _ := strconv.ParseInt(strings.TrimPrefix(data, "program_delete_"), 10, 64)
		h.deleteProgram(chatID, programID)
	}
}

func (h *Handler) programManagement(chatID int64) {
	res, err := h.findAllProgramsByUserUC.Execute(chatID)
	if err != nil {
		if errors.Is(err, programusecases.NoProgramsErr) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.NoProgramsFound)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.findAllProgramsByUserUC.Name())
		return
	}
	h.presenter.ShowProgramManageDialog(chatID, res)
}

func (h *Handler) createProgram(chatID int64) {
	if err := h.createProgramUC.Execute(chatID); err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.createProgramUC.Name())
		return
	}
	programsResult, err := h.findAllProgramsByUserUC.Execute(chatID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findAllProgramsByUserUC.Name())
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.SuccessfullyCreatedProgram)
	h.presenter.ShowProgramManageDialog(chatID, programsResult)
}

func (h *Handler) activateProgram(chatID int64, programID int64) {
	if err := h.activateProgramUC.Execute(chatID, programID); err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.activateProgramUC.Name())
		return
	}
	programsResult, err := h.findAllProgramsByUserUC.Execute(chatID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findAllProgramsByUserUC.Name())
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.SuccessfullyChangedProgram)
	h.presenter.ShowProgramManageDialog(chatID, programsResult)
}

func (h *Handler) deleteProgram(chatID int64, programID int64) {
	err := h.deleteProgramUC.Execute(chatID, programID)
	if err != nil {
		if errors.Is(err, programusecases.CannotDeleteCurrentProgramErr) {
			h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.CannotDeleteCurrentProgram)
			return
		}
		h.commonPresenter.HandleInternalError(err, chatID, h.deleteProgramUC.Name())
		return
	}
	programsResult, err := h.findAllProgramsByUserUC.Execute(chatID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.findAllProgramsByUserUC.Name())
		return
	}
	h.commonPresenter.SendSimpleHtmlMessage(chatID, messages.SuccessfullyDeletedProgram)
	h.presenter.ShowProgramManageDialog(chatID, programsResult)
}

func (h *Handler) editProgram(chatID int64, programID int64) {
	res, err := h.getProgramUC.Execute(programID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getProgramUC.Name())
		return
	}
	h.presenter.ShowEditDialog(chatID, res)
}

func (h *Handler) confirmDeleteProgram(chatID int64, programID int64) {
	res, err := h.getProgramUC.Execute(programID)
	if err != nil {
		h.commonPresenter.HandleInternalError(err, chatID, h.getProgramUC.Name())
		return
	}
	h.presenter.ConfirmDeleteDialog(chatID, res)
}
