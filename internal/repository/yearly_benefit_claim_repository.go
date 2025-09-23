package repository

import (
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
	return &YearlyBenefitClaimRepository{Log: log}
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
	if err := base.Find(&yearlyBenefitClaims).Error; err != nil {
		r.Log.WithField("repository", "YearlyBenefitClaimRepository").Errorf("get all yearly benefit claims failed: %v", err)
		return nil, 0, err
	}
	return yearlyBenefitClaims, count, nil
}
