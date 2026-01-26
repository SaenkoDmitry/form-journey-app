package users

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/programs"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CreateUseCase struct {
	usersRepo    users.Repo
	programsRepo programs.Repo
}

func NewCreateUseCase(usersRepo users.Repo, programsRepo programs.Repo) *CreateUseCase {
	return &CreateUseCase{
		usersRepo:    usersRepo,
		programsRepo: programsRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Создать пользователя"
}

func (uc *CreateUseCase) Execute(chatID int64, from *tgbotapi.User) (*models.User, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		if errors.Is(err, users.NotFoundUserErr) {
			createdUser, createErr := uc.usersRepo.Create(chatID, from)
			if createErr != nil {
				return nil, createErr
			}

			// создаем дефолтную программу
			program, createErr := uc.programsRepo.Create(createdUser.ID, "#1 стартовая")
			if createErr != nil {
				return nil, createErr
			}

			// прикрепляем программу к юзеру и сохраняем
			createdUser.ActiveProgramID = &program.ID
			err = uc.usersRepo.Save(createdUser)
			if err != nil {
				return nil, err
			}
			return createdUser, nil
		}
		return nil, err
	}

	return user, nil
}
