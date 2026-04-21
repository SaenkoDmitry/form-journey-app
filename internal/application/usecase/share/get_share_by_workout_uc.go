package share

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"github.com/SaenkoDmitry/training-tg-bot/internal/repository/share"
)

type GetShareByWorkoutUC struct {
	shareRepo share.Repo
}

func NewGetShareByWorkoutUC(shareRepo share.Repo) *GetShareByWorkoutUC {
	return &GetShareByWorkoutUC{shareRepo: shareRepo}
}

func (uc *GetShareByWorkoutUC) Execute(workoutID int64) (*models.WorkoutShare, error) {
	shareModel, err := uc.shareRepo.GetByWorkoutID(workoutID)
	if err != nil {
		return nil, err
	}
	return &shareModel, nil
}
