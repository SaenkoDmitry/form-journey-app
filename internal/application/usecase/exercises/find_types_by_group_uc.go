package exercises

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/exercisetypes"
)

type FindTypesByGroupUseCase struct {
	exerciseTypesRepo exercisetypes.Repo
}

func NewFindTypesByGroupUseCase(exerciseTypesRepo exercisetypes.Repo) *FindTypesByGroupUseCase {
	return &FindTypesByGroupUseCase{
		exerciseTypesRepo: exerciseTypesRepo,
	}
}

func (uc *FindTypesByGroupUseCase) Name() string {
	return "Найти упражнения по группе"
}

func (uc *FindTypesByGroupUseCase) Execute(exerciseGroupCode string) (*dto.FindTypesByGroup, error) {
	exerciseTypes, err := uc.exerciseTypesRepo.GetAllByGroup(exerciseGroupCode)
	if err != nil {
		return nil, err
	}

	return &dto.FindTypesByGroup{
		ExerciseTypes: exerciseTypes,
	}, nil
}
