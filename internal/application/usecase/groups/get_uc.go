package groups

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
)

type GetUseCase struct {
	exerciseGroupTypesRepo exercisegrouptypes.Repo
}

func NewGetUseCase(
	exerciseGroupTypesRepo exercisegrouptypes.Repo,
) *GetUseCase {
	return &GetUseCase{
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
	}
}

func (uc *GetUseCase) Name() string {
	return "Получить информацию о группе упражнений"
}

func (uc *GetUseCase) Execute(exerciseGroupCode string) (*dto.Group, error) {
	res, err := uc.exerciseGroupTypesRepo.Get(exerciseGroupCode)
	if err != nil {
		return nil, err
	}
	return &dto.Group{
		Code: res.Code,
		Name: res.Name,
	}, nil
}
