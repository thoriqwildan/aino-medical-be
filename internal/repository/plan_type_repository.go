package repository

import (
	"errors"

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

func (ptr *PlanTypeRepository) FindOrCreate(db *gorm.DB, planType *entity.PlanType) (*entity.PlanType, error) {
	var result entity.PlanType
	errTake := db.Model(entity.PlanType{}).Where("name = ?", planType.Name).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		ptr.Log.Error("Error when find plan type in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil 
	}

	errCreate := db.Create(planType).Error
	if errCreate != nil {
		ptr.Log.Error("Error when create plan type in method find or create")
		return nil, errCreate
	}

	return planType, nil
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
	if request.Limit > 0 {
		baseQuery = baseQuery.Limit(request.Limit)
	}
	if request.Page > 0 {
		baseQuery = baseQuery.Offset((request.Page - 1) * request.Limit)
	}
	err := baseQuery.
		Find(&planTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return planTypes, total, nil
}
