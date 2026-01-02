package sessions

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(session *models.WorkoutSession) error
	GetByWorkoutID(sessionID int64) (models.WorkoutSession, error)
	Save(session *models.WorkoutSession) error
	UpdateIsActive(workoutID int64, isActive bool) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Create(session *models.WorkoutSession) error {
	u.db.Create(&session)
	return nil
}

func (u *repoImpl) GetByWorkoutID(workoutID int64) (models.WorkoutSession, error) {
	var session models.WorkoutSession
	u.db.Where("workout_day_id = ? AND is_active = ?", workoutID, true).
		Order("started_at DESC").
		First(&session)
	return session, nil
}

func (u *repoImpl) Save(session *models.WorkoutSession) error {
	u.db.Save(&session)
	return nil
}

func (u *repoImpl) UpdateIsActive(workoutID int64, isActive bool) error {
	u.db.Model(&models.WorkoutSession{}).
		Where("workout_day_id = ? AND is_active = ?", workoutID, true).
		Update("is_active", isActive)
	return nil
}
