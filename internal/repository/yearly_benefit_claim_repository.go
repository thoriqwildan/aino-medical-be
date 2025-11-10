package repository

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type YearlyBenefitClaimRepository struct {
	Repository[entity.YearlyBenefitClaim]
	Log *logrus.Logger
}

func NewYearlyBenefitClaimRepository(log *logrus.Logger) *YearlyBenefitClaimRepository {
	return &YearlyBenefitClaimRepository{Log: log, Repository: Repository[entity.YearlyBenefitClaim]{}}
}

func (r YearlyBenefitClaimRepository) GetAll(db *gorm.DB, request *model.YearlyBenefitClaimFilter) ([]*entity.YearlyBenefitClaim, int64, error) {
	var yearlyBenefitClaims []*entity.YearlyBenefitClaim
	base := db.Model(&entity.YearlyBenefitClaim{})
	if request != nil {
		if request.Code != "" {
			base = base.Where("code LIKE ?", "%"+request.Code+"%")
		}
		if request.Limit != 0 {
			base = base.Limit(request.Limit)
		}
		if request.Page != 0 {
			base = base.Offset((request.Page - 1) * request.Limit)
		}
	}
	var count int64
	if err := base.Count(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("error when counting result find yearly benefit claims: %v", err)
	}
	if err := base.Preload("Benefits").Find(&yearlyBenefitClaims).Error; err != nil {
		r.Log.WithField("repository", "YearlyBenefitClaimRepository").Errorf("get all yearly benefit claims failed: %v", err)
		return nil, 0, err
	}
	return yearlyBenefitClaims, count, nil
}

func (r *YearlyBenefitClaimRepository) FindOrCreate(db *gorm.DB, yearlyBenefitClaim *entity.YearlyBenefitClaim) (*entity.YearlyBenefitClaim, error) {
	var result entity.YearlyBenefitClaim
	errTake := db.Model(entity.YearlyBenefitClaim{}).Where("name = ?", yearlyBenefitClaim.Code).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		r.Log.Error("Error when find yearly benefit claim in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil 
	}

	errCreate := db.Create(yearlyBenefitClaim).Error
	if errCreate != nil {
		r.Log.Error("Error when create yearly benefit claim in method find or create")
		return nil, errCreate
	}

	return yearlyBenefitClaim, nil
}