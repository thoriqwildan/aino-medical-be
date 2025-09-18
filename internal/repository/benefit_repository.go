package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
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

func (br *BenefitRepository) GetByCode(db *gorm.DB, code string) error {
	return db.Where("code = ?", code).First(&entity.Benefit{}).Error
}

func (br *BenefitRepository) GetById(db *gorm.DB, id uint, benefit *entity.Benefit) error {
	return db.Where("id = ?", id).Preload("PlanType").First(benefit).Error
}

func (br *BenefitRepository) SearchBenefits(db *gorm.DB, request *model.SearchPagingQuery) ([]entity.Benefit, int64, error) {
	var benefits []entity.Benefit
	var total int64

	baseQuery := db.Model(&entity.Benefit{})

	if request.SearchValue != "" {
		baseQuery.Where("name LIKE ?", "%"+request.SearchValue+"%")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Preload("PlanType").
		Preload("YearlyBenefitClaim").
		Find(&benefits).Error
	if err != nil {
		return nil, 0, err
	}

	return benefits, total, nil
}
