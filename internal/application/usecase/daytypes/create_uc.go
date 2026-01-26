package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"time"
)

type CreateUseCase struct {
	dayTypesRepo daytypes.Repo
}

func NewCreateUseCase(
	dayTypesRepo daytypes.Repo,
) *CreateUseCase {
	return &CreateUseCase{
		dayTypesRepo: dayTypesRepo,
	}
}

func (uc *CreateUseCase) Name() string {
	return "Создать новый тренировочный день"
}

func (uc *CreateUseCase) Execute(programID int64, name string) (int64, error) {
	dayType, createErr := uc.dayTypesRepo.Create(&models.WorkoutDayType{
		WorkoutProgramID: programID,
		Name:             name,
		CreatedAt:        time.Now(),
	})
	if createErr != nil {
		return 0, createErr
	}

	return dayType.ID, nil
}
