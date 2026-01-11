package programs

import (
	"time"

	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(userID int64, name string) (*models.WorkoutProgram, error)
	Save(program *models.WorkoutProgram) error
	Get(programID int64) (models.WorkoutProgram, error)
	Delete(program *models.WorkoutProgram) error
	FindAll(userID int64) ([]models.WorkoutProgram, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) Create(userID int64, name string) (*models.WorkoutProgram, error) {
	newProgram := &models.WorkoutProgram{
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	err := u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&newProgram).Error
	})
	return newProgram, err
}

func (u *repoImpl) Delete(program *models.WorkoutProgram) error {
	return u.db.Delete(program).Error
}

func (u *repoImpl) Save(program *models.WorkoutProgram) error {
	return u.db.Save(program).Error
}

func (u *repoImpl) Get(programID int64) (program models.WorkoutProgram, err error) {
	tx := u.db.
		Preload("DayTypes", func(db *gorm.DB) *gorm.DB { return db.Order("created_at ASC") }).
		//Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		//Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		First(&program, programID)
	if tx.Error != nil {
		return models.WorkoutProgram{}, tx.Error
	}

	return program, nil
}

func (u *repoImpl) FindAll(userID int64) (programs []models.WorkoutProgram, err error) {
	tx := u.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		//Preload("Exercises.Sets", func(db *gorm.DB) *gorm.DB { return db.Order("sets.index ASC") }).
		//Preload("Exercises", func(db *gorm.DB) *gorm.DB { return db.Order("exercises.index ASC") }).
		Find(&programs)

	if tx.Error != nil {
		return []models.WorkoutProgram{}, tx.Error
	}

	return programs, nil
}
