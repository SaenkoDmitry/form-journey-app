package groups

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisegrouptypes"
)

type GetAllUseCase struct {
	exerciseGroupTypesRepo exercisegrouptypes.Repo
}

func NewGetAllUseCase(exerciseGroupTypesRepo exercisegrouptypes.Repo) *GetAllUseCase {
	return &GetAllUseCase{
		exerciseGroupTypesRepo: exerciseGroupTypesRepo,
	}
}

func (uc *GetAllUseCase) Name() string {
	return "Список групп упражнений"
}

func (uc *GetAllUseCase) Execute() (*dto.ExerciseGroupTypeList, error) {
	groups, err := uc.exerciseGroupTypesRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return &dto.ExerciseGroupTypeList{
		Groups: groups,
	}, nil
}
