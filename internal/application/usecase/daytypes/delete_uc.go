package daytypes

import (
	"errors"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/daytypes"
	"strings"
)

type DeleteUseCase struct {
	dayTypesRepo daytypes.Repo
}

func NewDeleteUseCase(
	dayTypesRepo daytypes.Repo,
) *DeleteUseCase {
	return &DeleteUseCase{
		dayTypesRepo: dayTypesRepo,
	}
}

func (uc *DeleteUseCase) Name() string {
	return "Удалить тренировочный день"
}

var (
	CannotDeleteAlreadyUsedDay = errors.New("cannot delete day already used")
)

func (uc *DeleteUseCase) Execute(dayTypeID int64) error {
	deleteErr := uc.dayTypesRepo.Delete(dayTypeID)
	if deleteErr != nil {
		if strings.Contains(deleteErr.Error(), "update or delete on table \"workout_day_types\" violates foreign key constraint") {
			return CannotDeleteAlreadyUsedDay
		}
		return deleteErr
	}
	return nil
}
