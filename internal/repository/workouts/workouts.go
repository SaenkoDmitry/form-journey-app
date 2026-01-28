package workouts

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(workoutDay *models.WorkoutDay) error
	Delete(workout *models.WorkoutDay) error
	Save(workout *models.WorkoutDay) error
	Get(workoutID int64) (models.WorkoutDay, error)
	Count(userID int64) (count int64, err error)
	FindAll(userID int64) ([]models.WorkoutDay, error)
	Find(userID int64, offset, limit int) ([]models.WorkoutDay, error)
	FindPreviousByType(userID int64, dayTypeID int64) (models.WorkoutDay, error)
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
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(workoutDay).Error
	})
}

func (u *repoImpl) Delete(workout *models.WorkoutDay) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(workout).Error
	})
}

func (u *repoImpl) Save(workout *models.WorkoutDay) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(workout).Error
	})
}

func (u *repoImpl) Get(workoutID int64) (workoutDay models.WorkoutDay, err error) {
	err = u.db.Transaction(func(tx *gorm.DB) error {
		return tx.
			Preload("WorkoutDayType").
			Preload("Exercises.WorkoutDay").
			Preload("Exercises.ExerciseType").
			Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
			Preload("Exercises.Sets.Exercise.ExerciseType").
			Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
			First(&workoutDay, workoutID).Error
	})
	return workoutDay, err
}

func (u *repoImpl) Count(userID int64) (count int64, err error) {
	err = u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.WorkoutDay{}).Where("user_id = ?", userID).Count(&count).Error
	})
	return count, err
}

func (u *repoImpl) FindAll(userID int64) (workouts []models.WorkoutDay, err error) {
	err = u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", userID).
			Order("started_at DESC").
			Preload("WorkoutDayType").
			Preload("Exercises.ExerciseType").
			Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
			Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
			Find(&workouts).Error
	})
	return workouts, err
}

func (u *repoImpl) Find(userID int64, offset, limit int) (workouts []models.WorkoutDay, err error) {
	err = u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", userID).
			Order("started_at DESC").
			Preload("WorkoutDayType").
			Preload("Exercises.ExerciseType").
			Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
			Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
			Offset(offset).
			Limit(limit).
			Find(&workouts).Error
	})

	return workouts, err
}

func (u *repoImpl) FindPreviousByType(userID int64, dayTypeID int64) (workout models.WorkoutDay, err error) {
	err = u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ? AND workout_day_type_id = ? AND completed = ?", userID, dayTypeID, true).
			Order("started_at DESC").
			Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
			Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
			First(&workout).Error
	})
	return workout, err
}
