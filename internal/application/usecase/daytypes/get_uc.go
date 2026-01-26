package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
)

type GetUseCase struct {
	dayTypesRepo daytypes.Repo
}

func NewGetUseCase(
	dayTypesRepo daytypes.Repo,
) *GetUseCase {
	return &GetUseCase{
		dayTypesRepo: dayTypesRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Загрузить тренировочный день"
}

func (uc *GetUseCase) Execute(dayTypeID int64) (*models.WorkoutDayType, error) {
	dayType, err := uc.dayTypesRepo.Get(dayTypeID)
	if err != nil {
		return nil, err
	}

	return &dayType, nil
}
