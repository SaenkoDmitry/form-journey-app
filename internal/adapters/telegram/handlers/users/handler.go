package users

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/adapters/telegram/common"
	userusecases "github.com/SaenkoDmitry/training-tg-bot/internal/application/usecase/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	presenter       *Presenter
	commonPresenter *common.Presenter

	createUserUC *userusecases.CreateUseCase
}

func NewHandler(bot *tgbotapi.BotAPI, createUserUC *userusecases.CreateUseCase) *Handler {
	return &Handler{
		presenter:       NewPresenter(bot),
		commonPresenter: common.NewPresenter(bot),
		createUserUC:    createUserUC,
	}
}

func (h *Handler) RouteCallback(chatID int64, data string) {
	switch {

	}
}
