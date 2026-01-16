package sets

import (
	"github.com/SaenkoDmitry/training-tg-bot/internal/models"
	"gorm.io/gorm"
)

type Repo interface {
	Delete(id int64) error
	DeleteAllBy(exerciseID int64) error
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

func (u *repoImpl) Delete(id int64) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&models.Set{}).Error
	})
}

func (u *repoImpl) DeleteAllBy(exerciseID int64) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("exercise_id = ?", exerciseID).Delete(&models.Set{}).Error
	})
}

func (u *repoImpl) Save(set *models.Set) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&set).Error
	})
}
