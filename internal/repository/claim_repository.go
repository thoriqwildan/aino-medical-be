package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type ClaimRepository struct {
	Repository[entity.Claim]
	Log *logrus.Logger
}

func NewClaimRepository(log *logrus.Logger) *ClaimRepository {
	return &ClaimRepository{
		Log: log,
	}
}

func(r *ClaimRepository) GetBenefitByCode(db *gorm.DB, benefit *entity.Benefit, code string) error {
	return db.Where("code = ?", code).
				Preload("PlanType").
				Preload("LimitationType").
				First(benefit).Error
}

func (r *ClaimRepository) GetByID(db *gorm.DB, claim *entity.Claim, id any) error {
	return db.Where("id = ?", id).
				Preload("Patient").
				Preload("PatientBenefit.Benefit").
				Preload("PatientBenefit.Benefit.PlanType").
				Preload("PatientBenefit.Benefit.LimitationType").
				Preload("TransactionType").
				First(claim).Error
}

func (r *ClaimRepository) GetPatientByID(db *gorm.DB, patient *entity.Patient, id any) error {
	return db.Where("id = ?", id).
				Preload("Employee").
				Preload("FamilyMember").
				First(patient).Error
}

func (r *ClaimRepository) GetPatients(db *gorm.DB, request *model.PagingQuery) ([]entity.Patient, int64, error) {
	var patients []entity.Patient
	var total int64

	baseQuery := db.Model(&entity.Patient{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Preload("PlanType").
		Find(&patients).Error
	if err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

func (r *ClaimRepository) GetBenefits(db *gorm.DB, request *model.PagingQuery, planTypeID uint) ([]entity.Benefit, int64, error) {
	var benefits []entity.Benefit
	var total int64

	baseQuery := db.Model(&entity.Benefit{})

	if err := baseQuery.Where("plan_type_id = ?", planTypeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Where("plan_type_id = ?", planTypeID).
		Preload("PlanType").
		Find(&benefits).Error
	if err != nil {
		return nil, 0, err
	}

	return benefits, total, nil
}
