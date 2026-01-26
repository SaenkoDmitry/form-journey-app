package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
)

type UpdateUseCase struct {
	dayTypesRepo daytypes.Repo
}

func NewUpdateUseCase(
	dayTypesRepo daytypes.Repo,
) *UpdateUseCase {
	return &UpdateUseCase{
		dayTypesRepo: dayTypesRepo,
	}
}

func (uc *UpdateUseCase) Name() string {
	return "Обновить день"
}

func (uc *UpdateUseCase) Execute(dayType *models.WorkoutDayType) error {
	err := uc.dayTypesRepo.Save(dayType)
	if err != nil {
		return err
	}
	return nil
}
