package sets

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	GetCompleted(exerciseID int64) int64
	Delete(exerciseID int64) error
	Save(set *models.Set) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (u *repoImpl) GetCompleted(exerciseID int64) int64 {
	var completedSets int64
	u.db.Model(&models.Set{}).Where("exercise_id = ? AND completed = ?", exerciseID, true).Count(&completedSets)
	return completedSets
}

func (u *repoImpl) Delete(exerciseID int64) error {
	u.db.Where("exercise_id = ?", exerciseID).Delete(&models.Set{})
	return nil
}

func (u *repoImpl) Save(set *models.Set) error {
	u.db.Save(&set)
	return nil
}
