package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type LimitationTypeRepository struct {
	Repository[entity.LimitationType]
	Log *logrus.Logger
}

func NewLimitationTypeRepository(log *logrus.Logger) *LimitationTypeRepository {
	return &LimitationTypeRepository{
		Log: log,
	}
}

func (r *LimitationTypeRepository) GetByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.LimitationType{}).Error
}

func (r *LimitationTypeRepository) SearchLimitationTypes(db *gorm.DB, request *model.PagingQuery) ([]entity.LimitationType, int64, error) {
	var limitationTypes []entity.LimitationType
	var total int64

	baseQuery := db.Model(&entity.LimitationType{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Find(&limitationTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return limitationTypes, total, nil
}