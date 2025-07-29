package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

type BenefitRepository struct {
	Repository[entity.Benefit]
	Log *logrus.Logger
}

func NewBenefitRepository(log *logrus.Logger) *BenefitRepository {
	return &BenefitRepository{
		Log: log,
	}
}

func (br *BenefitRepository) GetByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.Benefit{}).Error
}