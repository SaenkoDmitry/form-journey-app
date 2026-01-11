package sets

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
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

func (u *repoImpl) Delete(exerciseID int64) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("exercise_id = ?", exerciseID).Delete(&models.Set{}).Error
	})
}

func (u *repoImpl) Save(set *models.Set) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&set).Error
	})
}
