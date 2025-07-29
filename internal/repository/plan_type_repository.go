package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type PlanTypeRepository struct {
	Repository[entity.PlanType]
	Log *logrus.Logger
}

func NewPlanTypeRepository(log *logrus.Logger) *PlanTypeRepository {
	return &PlanTypeRepository{
		Log: log,
	}
}

func (ptr *PlanTypeRepository) FindByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.PlanType{}).Error
}

func (ptr *PlanTypeRepository) SearchPlanTypes(db *gorm.DB, request *model.PagingQuery) ([]entity.PlanType, int64, error) {
	var planTypes []entity.PlanType
	var total int64

	baseQuery := db.Model(&entity.PlanType{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Find(&planTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return planTypes, total, nil
}