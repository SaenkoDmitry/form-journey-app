package measurements

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/application/dto"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/measurements"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/users"
	"github.com/SaenkoDmitry/training-tg-bot/internal/utils"
)

type FindAllByUserUseCase struct {
	measurementsRepo measurements.Repo
	usersRepo        users.Repo
}

func NewFindAllByUserUseCase(
	measurementsRepo measurements.Repo,
	usersRepo users.Repo,
) *FindAllByUserUseCase {
	return &FindAllByUserUseCase{
		measurementsRepo: measurementsRepo,
		usersRepo:        usersRepo,
	}
}

func (uc *FindAllByUserUseCase) Name() string {
	return "Показать измерения"
}

func (uc *FindAllByUserUseCase) Execute(chatID int64, limit, offset int) (*dto.FindWithOffsetLimitMeasurement, error) {
	user, err := uc.usersRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	measurementObjs, err := uc.measurementsRepo.FindAllLimitOffset(user.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	count, _ := uc.measurementsRepo.Count(user.ID)

	result := make([]dto.Measurement, 0, len(measurementObjs))
	for _, m := range measurementObjs {
		result = append(result, dto.Measurement{
			CreatedAt: utils.FormatDate(m.CreatedAt),
			Shoulders: utils.FormatCentimeters(m.Shoulders),
			Chest:     utils.FormatCentimeters(m.Chest),
			HandLeft:  utils.FormatCentimeters(m.HandLeft),
			HandRight: utils.FormatCentimeters(m.HandRight),
			Waist:     utils.FormatCentimeters(m.Waist),
			Buttocks:  utils.FormatCentimeters(m.Buttocks),
			HipLeft:   utils.FormatCentimeters(m.HipLeft),
			HipRight:  utils.FormatCentimeters(m.HipRight),
			CalfLeft:  utils.FormatCentimeters(m.CalfLeft),
			CalfRight: utils.FormatCentimeters(m.CalfRight),
			Weight:    utils.FormatKilograms(m.Weight),
		})
	}

	return &dto.FindWithOffsetLimitMeasurement{
		Items: result,
		Count: int(count),
	}, nil
}
