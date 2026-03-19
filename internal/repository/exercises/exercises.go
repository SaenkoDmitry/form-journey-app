package exercises

import (
	"fmt"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Get(exerciseID int64) (models.Exercise, error)
	FindAllByWorkoutID(workoutDayID int64) ([]models.Exercise, error)
	FindAllByUserIDAndExTypeID(userID int64, exerciseTypeID int64, offset, limit int) ([]models.Exercise, error)
	CountByUserIDAndExTypeID(userID int64, exerciseTypeID int64) (int64, error)
	DeleteByWorkout(workoutID int64) error
	Delete(exerciseID int64) error
	CreateBatch(exercises []models.Exercise) error
	Save(exercise *models.Exercise) error
	FindPreviousByType(exerciseTypeID, activeProgramID int64) (models.Exercise, error)
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
	u.db.Preload("WorkoutDay").Preload("ExerciseType").Preload("Sets", func(db *gorm.DB) *gorm.DB {
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

func (u *repoImpl) FindAllByUserIDAndExTypeID(userID int64, exerciseTypeID int64, offset, limit int) ([]models.Exercise, error) {
	var exercises []models.Exercise

	err := u.db.
		Joins("JOIN workout_days ON workout_days.id = exercises.workout_day_id").
		Where("workout_days.user_id = ? AND exercise_type_id = ?", userID, exerciseTypeID).
		Preload("ExerciseType").
		Preload("WorkoutDay").
		Preload("Sets.Exercise.ExerciseType").
		Preload("Sets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sets.index ASC")
		}).
		Order("workout_days.started_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&exercises).Error

	return exercises, err
}

func (u *repoImpl) CountByUserIDAndExTypeID(userID int64, exerciseTypeID int64) (int64, error) {
	var count int64

	err := u.db.
		Joins("JOIN workout_days ON workout_days.id = exercises.workout_day_id").
		Where("workout_days.user_id = ? AND exercise_type_id = ?", userID, exerciseTypeID).
		Table("exercises").
		Count(&count).Error

	return count, err
}

func (u *repoImpl) FindPreviousByType(exerciseTypeID, activeProgramID int64) (models.Exercise, error) {
	var exercise models.Exercise
	err := u.db.Joins("JOIN sets ON sets.exercise_id = exercises.id"+
		" JOIN workout_days wd ON wd.id = workout_day_id JOIN workout_day_types wdt ON wdt.id = wd.workout_day_type_id").
		Where("exercise_type_id = ? AND sets.completed = true AND wd.completed = true AND wdt.workout_program_id = ?", exerciseTypeID, activeProgramID).
		Preload("ExerciseType").
		Preload("Sets", func(db *gorm.DB) *gorm.DB {
			return db.Order("sets.index ASC")
		}).
		Order("sets.completed_at DESC").
		First(&exercise).Error
	return exercise, err
}
