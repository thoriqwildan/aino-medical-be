package repository

import (
	"time"

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
				Preload("PlanType").
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
		Preload("Employee").
		Preload("FamilyMember").
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
		Preload("LimitationType").
		Preload("PlanType").
		Find(&benefits).Error
	if err != nil {
		return nil, 0, err
	}

	return benefits, total, nil
}

func (r *ClaimRepository) FindAllWithQuery(db *gorm.DB, query *model.ClaimFilterQuery) ([]entity.Claim, int64, error) {
    var claims []entity.Claim
    var total int64

    // Hitung total data (tanpa limit dan offset)
    countDB := db.Model(&entity.Claim{})
    countDB = r.applyFilters(countDB, query)
    countDB.Count(&total)

    // Terapkan filter, preload, limit, dan offset
    queryDB := db.Model(&entity.Claim{})
    queryDB = r.applyFilters(queryDB, query)
    
    // Terapkan Preload yang Anda butuhkan
    queryDB = queryDB.
        Preload("Patient").
        Preload("Employee").
        Preload("PatientBenefit.Benefit").
        Preload("TransactionType")

    // Terapkan pagination
    offset := (query.Page - 1) * query.Limit
    if offset < 0 {
        offset = 0
    }
    
    if query.Limit > 0 {
        queryDB = queryDB.Limit(query.Limit).Offset(offset)
    }

    err := queryDB.Find(&claims).Error
    if err != nil {
        return nil, 0, err
    }

    return claims, total, nil
}

func (r *ClaimRepository) applyFilters(db *gorm.DB, query *model.ClaimFilterQuery) *gorm.DB {
    if query.TransactionStatus != "" {
        db = db.Where("transaction_status = ?", query.TransactionStatus)
    }
    
    if query.ClaimStatus != "" {
        db = db.Where("claim_status = ?", query.ClaimStatus)
    }

    if query.SLAStatus != "" {
        db = db.Where("SLA = ?", query.SLAStatus)
    }
    
    if query.TransactionType != "" {
        db = db.Joins("TransactionType").Where("TransactionType.name = ?", query.TransactionType)
    }

    if query.Department != "" {
        db = db.Joins("Employee.Department").Where("Department.name = ?", query.Department)
    }

    if query.DateFrom != "" {
        if date, err := time.Parse("2006-01-02", query.DateFrom); err == nil {
            db = db.Where("transaction_date >= ?", date)
        }
    }
    if query.DateTo != "" {
        if date, err := time.Parse("2006-01-02", query.DateTo); err == nil {
            db = db.Where("transaction_date <= ?", date)
        }
    }

    return db
}
