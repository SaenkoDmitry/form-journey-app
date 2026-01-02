package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(workoutDay *models.WorkoutDay) error
	Get(workoutID int64) (models.WorkoutDay, error)
	Find(userID int64) ([]models.WorkoutDay, error)
	Delete(workout *models.WorkoutDay) error
	Save(workout *models.WorkoutDay) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Create(workoutDay *models.WorkoutDay) error {
	u.db.Create(&workoutDay)
	return nil
}

func (u *repoImpl) Get(workoutID int64) (models.WorkoutDay, error) {
	var workoutDay models.WorkoutDay
	u.db.Preload("Exercises.Sets").First(&workoutDay, workoutID)
	return workoutDay, nil
}

func (u *repoImpl) Find(userID int64) ([]models.WorkoutDay, error) {
	var workouts []models.WorkoutDay
	u.db.Where("user_id = ?", userID).
		Order("started_at DESC").
		Find(&workouts)
	return workouts, nil
}

func (u *repoImpl) Delete(workout *models.WorkoutDay) error {
	u.db.Delete(workout)
	return nil
}

func (u *repoImpl) Save(workout *models.WorkoutDay) error {
	u.db.Save(workout)
	return nil
}
