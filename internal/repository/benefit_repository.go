package repository

import (
	"errors"

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

func (br *BenefitRepository) FindOrCreate(db *gorm.DB, benefit *entity.Benefit) (*entity.Benefit,error) {
	var result entity.Benefit
	errTake := db.Model(entity.Benefit{}).Where("code = ?", benefit.Code).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		br.Log.Error("Error when find benefit in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil 
	}

	errCreate := db.Create(benefit).Error
	if errCreate != nil {
		br.Log.Error("Error when create benefit in method find or create")
		return nil, errCreate
	}

	return benefit, nil
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

	baseQuery := db.Model(&entity.Benefit{}).
		Preload("YearlyBenefitClaim").
		Preload("PlanType")

	if request.SearchValue != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+request.SearchValue+"%")
	}

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
		Order("plan_type_id ASC").
		Find(&benefits).Error
	if err != nil {
		return nil, 0, err
	}

	return benefits, total, nil
}
