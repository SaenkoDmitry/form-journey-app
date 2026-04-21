package share

import (
	"gorm.io/gorm"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
)

type Repo interface {
	Create(share *models.WorkoutShare) error
	Get(token string) (share models.WorkoutShare, err error)
	GetByWorkoutID(workoutID int64) (share models.WorkoutShare, err error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Get(token string) (share models.WorkoutShare, err error) {
	tx := u.db.Where("token = ?", token).First(&share)
	if tx.Error != nil {
		return models.WorkoutShare{}, tx.Error
	}
	return share, nil
}

func (u *repoImpl) GetByWorkoutID(workoutID int64) (share models.WorkoutShare, err error) {
	tx := u.db.Where("workout_day_id = ?", workoutID).First(&share)
	if tx.Error != nil {
		return models.WorkoutShare{}, tx.Error
	}
	return share, nil
}

func (u *repoImpl) Create(share *models.WorkoutShare) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&share).Error
	})
}
