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

func (r *ClaimRepository) GetBenefitByCode(db *gorm.DB, benefit *entity.Benefit, code string) error {
	return db.Where("code = ?", code).
		Preload("PlanType").
		Preload("LimitationType").
		First(benefit).Error
}

func (r *ClaimRepository) GetByID(db *gorm.DB, claim *entity.Claim, id any) error {
	return db.Where("id = ?", id).
		Preload("Patient").
		Preload("Employee").
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
		Preload("FamilyMember.Employee").
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
		Offset((request.Page-1)*request.Limit).
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

func (r *ClaimRepository) GetBenefitsWithPlafond(db *gorm.DB, request *model.PagingQuery, planTypeID uint, patientID uint) ([]entity.Benefit, map[uint]float64, int64, error) {
	var benefits []entity.Benefit
	var total int64

	// 1. Ambil semua benefit berdasarkan plan_type_id
	baseQuery := db.Model(&entity.Benefit{}).Where("plan_type_id = ?", planTypeID)

	// Hitung total data
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	// Terapkan pagination dan preload
	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Preload("LimitationType").
		Preload("PlanType").
		Find(&benefits).Error

	if err != nil {
		return nil, nil, 0, err
	}

	// 2. Ambil semua patient_benefit untuk pasien tersebut
	var patientBenefits []entity.PatientBenefit
	benefitIDs := make([]uint, len(benefits))
	for i, b := range benefits {
		benefitIDs[i] = b.ID
	}

	db.Where("patient_id = ? AND benefit_id IN ?", patientID, benefitIDs).Find(&patientBenefits)

	// 3. Buat map untuk memudahkan pencarian remaining_plafond
	remainingPlafondMap := make(map[uint]float64)
	for _, pb := range patientBenefits {
		if pb.RemainingPlafond != nil {
			remainingPlafondMap[pb.BenefitID] = *pb.RemainingPlafond
		} else {
			remainingPlafondMap[pb.BenefitID] = 0
		}
	}

	return benefits, remainingPlafondMap, total, nil
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
		Preload("Patient.PlanType").
		Preload("Employee").
		Preload("Employee.PlanType").
		Preload("Employee.Department").
		Preload("Employee.FamilyMembers").
		Preload("Employee.FamilyMembers.Employee").
		Preload("PatientBenefit.Benefit").
		Preload("PatientBenefit.Benefit.PlanType").
		Preload("PatientBenefit.Benefit.LimitationType").
		Preload("TransactionType").
		Order("COALESCE(claims.updated_at, claims.created_at) DESC")

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
	// Track which joins we need to avoid duplicates
	needsEmployeeJoin := query.SearchValue != "" || query.Department != ""
	needsPatientJoin := query.SearchValue != "" || query.RelationshipType != ""

	// Apply joins first
	if needsEmployeeJoin {
		db = db.Joins("Employee")
	}
	if needsPatientJoin {
		db = db.Joins("Patient")
	}

	// Apply filters
	if query.SearchValue != "" {
		likeValue := "%" + query.SearchValue + "%"
		db = db.Where(
			db.Or("Employee.name LIKE ?", likeValue).
				Or("Patient.name LIKE ?", likeValue),
		)
	}

	if query.RelationshipType != "" {
		db = db.Joins("LEFT JOIN family_members ON Patient.family_member_id = family_members.id").
			Where("family_members.relationship_type = ?", query.RelationshipType)
	}

	if query.BenefitID != "" {
		db = db.Joins("LEFT JOIN patient_benefits ON patient_benefit_id = patient_benefits.id ").
			Where("patient_benefits.benefit_id = ?", query.BenefitID)
	}

	if query.TransactionStatus != "" {
		db = db.Where("claims.transaction_status = ?", query.TransactionStatus)
	}

	if query.ClaimStatus != "" {
		db = db.Where("claims.claim_status = ?", query.ClaimStatus)
	}

	if query.SLAStatus != "" {
		db = db.Where("claims.SLA = ?", query.SLAStatus)
	}

	if query.TransactionType != "" {
		db = db.Joins("TransactionType").Where("TransactionType.name = ?", query.TransactionType)
	}

	if query.Department != "" {
		db = db.Joins("JOIN departments ON Employee.department_id = departments.id").Where("departments.name = ?", query.Department)
	}

	if query.DateFrom != "" {
		if date, err := time.Parse("2006-01-02", query.DateFrom); err == nil {
			db = db.Where("claims.transaction_date >= ?", date)
		}
	}

	if query.DateTo != "" {
		if date, err := time.Parse("2006-01-02", query.DateTo); err == nil {
			db = db.Where("claims.transaction_date <= ?", date)
		}
	}

	return db
}
