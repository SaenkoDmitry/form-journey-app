package exercises

import (
	"fmt"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Get(exerciseID int64) (models.Exercise, error)
	FindAllByWorkoutID(workoutDayID int64) ([]models.Exercise, error)
	FindAllByUserID(userID int64) ([]models.Exercise, error)
	DeleteByWorkout(workoutID int64) error
	Delete(exerciseID int64) error
	CreateBatch(exercises []models.Exercise) error
	Save(exercise *models.Exercise) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Get(exerciseID int64) (models.Exercise, error) {
	var exercise models.Exercise
	u.db.Preload("ExerciseType").Preload("Sets", func(db *gorm.DB) *gorm.DB {
		return db.Order("sets.index ASC")
	}).First(&exercise, exerciseID)
	return exercise, nil
}

func (u *repoImpl) DeleteByWorkout(workoutID int64) error {
	u.db.Where("workout_day_id = ?", workoutID).Delete(&models.Exercise{})
	return nil
}

func (u *repoImpl) Delete(exerciseID int64) error {
	var exercise models.Exercise
	if err := u.db.Preload("Sets").First(&exercise, exerciseID).Error; err != nil {
		return err
	}

	// Удаляем с помощью Select
	return u.db.Select("Sets").Delete(&exercise).Error
}

func (u *repoImpl) Save(exercise *models.Exercise) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(exercise).Error
	})
}

func (u *repoImpl) CreateBatch(exercises []models.Exercise) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&exercises).Error
	})
}

func (u *repoImpl) FindAllByWorkoutID(workoutDayID int64) ([]models.Exercise, error) {
	var exercises []models.Exercise
	fmt.Println("FindAllByWorkoutID")

	u.db.Where("workout_day_id = ?", workoutDayID).
		Preload("Sets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sets.index ASC")
		}).
		Order("index ASC").Find(&exercises)

	return exercises, nil
}

func (u *repoImpl) FindAllByUserID(userID int64) ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := u.db.
		Joins("JOIN workout_days ON workout_days.id = exercises.workout_day_id").
		Where("workout_days.user_id = ?", userID).
		Preload("ExerciseType").
		Preload("Sets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sets.index ASC")
		}).
		Find(&exercises).Error

	return exercises, err
}
