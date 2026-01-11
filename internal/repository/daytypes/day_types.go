package daytypes

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(day *models.WorkoutDayType) (*models.WorkoutDayType, error)
	Delete(day *models.WorkoutDayType) error
	Save(day *models.WorkoutDayType) error
	Get(dayTypeID int64) (models.WorkoutDayType, error)
	FindAll(programID int64) ([]models.WorkoutDayType, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Create(day *models.WorkoutDayType) (*models.WorkoutDayType, error) {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&day).Error
	})
	return day, err
}

func (u *repoImpl) Delete(day *models.WorkoutDayType) error {
	u.db.Delete(day)
	return nil
}

func (u *repoImpl) Save(day *models.WorkoutDayType) error {
	u.db.Save(day)
	return nil
}

func (u *repoImpl) Get(dayTypeID int64) (day models.WorkoutDayType, err error) {
	u.db.First(&day, dayTypeID)
	return day, nil
}

func (u *repoImpl) FindAll(programID int64) (days []models.WorkoutDayType, err error) {
	u.db.
		Where("workout_program_id = ?", programID).
		Order("created_at DESC").
		Find(&days)

	return days, nil
}
