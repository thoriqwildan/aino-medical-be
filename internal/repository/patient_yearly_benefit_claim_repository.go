package repository

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type PatientYearlyBenefitClaimRepository struct {
	Repository[entity.PatientYearlyBenefitClaim]
	Log *logrus.Logger
}

func NewPatientYearlyBenefitClaimRepository(log *logrus.Logger) *PatientYearlyBenefitClaimRepository {

	return &PatientYearlyBenefitClaimRepository{Log: log}
}

func (r *PatientYearlyBenefitClaimRepository) GetAll(
	db *gorm.DB,
	request *model.SearchPagingQuery,
) ([]*entity.PatientYearlyBenefitClaim, int64, error) {

	var (
		rows  []*entity.PatientYearlyBenefitClaim
		total int64
	)

	base := db.Model(&entity.PatientYearlyBenefitClaim{}).
		Joins("JOIN patients p ON p.id = patient_yearly_benefit_claims.patient_id")

	filterScope := func(req *model.SearchPagingQuery) func(*gorm.DB) *gorm.DB {
		return func(tx *gorm.DB) *gorm.DB {
			if req == nil {
				return tx
			}
			if s := strings.TrimSpace(req.SearchValue); s != "" {
				tx = tx.Where("p.name LIKE ?", "%"+s+"%")
			}
			return tx
		}
	}

	if err := base.
		Scopes(filterScope(request)).
		Distinct("patient_yearly_benefit_claims.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := base.
		Scopes(filterScope(request)).
		Select("patient_yearly_benefit_claims.*").
		Preload("Patient").
		Preload("YearlyBenefitClaim").
		Order("p.name ASC")

	if request != nil && request.Page > 0 && request.Limit > 0 {
		q = q.Offset((request.Page - 1) * request.Limit).Limit(request.Limit)
	} else if request != nil && request.Limit > 0 {
		q = q.Limit(request.Limit)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *PatientYearlyBenefitClaimRepository) FindOrCreate(
	db *gorm.DB,
	patient entity.Patient,
	yearlyBenefitClaim entity.YearlyBenefitClaim,
) (*entity.PatientYearlyBenefitClaim, error) {
	var patientYearlyBenefitClaim entity.PatientYearlyBenefitClaim

	err := db.Model(entity.PatientYearlyBenefitClaim{}).Where("patient_id = ? AND yearly_benefit_claim_id = ?", patient.ID, yearlyBenefitClaim.ID).First(&patientYearlyBenefitClaim).Error

	if err == nil {
		r.Log.Printf("PatientYearlyBenefitClaim found for PatientID: %d, YearlyBenefitClaimID: %d", patient.ID, yearlyBenefitClaim.ID)
		return &patientYearlyBenefitClaim, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.Log.Printf("PatientYearlyBenefitClaim not foundfor PatientID: %d, YearlyBenefitClaimID: %d. Creating new record...", patient.ID, yearlyBenefitClaim.ID)

		newPatientYearlyBenefitClaim := entity.PatientYearlyBenefitClaim{
			PatientID:            patient.ID,
			YearlyBenefitClaimID: yearlyBenefitClaim.ID,
			YearlyClaimRemaining: helper.CalculateProratePlafond(yearlyBenefitClaim.YearlyClaim, patient.Employee.ProRate),
		}

		createErr := db.Create(&newPatientYearlyBenefitClaim).Error
		if createErr != nil {
			r.Log.Printf("Error creating PatientYearlyBenefitClaim for PatientID: %d, YearlyBenefitClaimID: %d: %v", patient.ID, yearlyBenefitClaim.ID, createErr)
			return nil, createErr
		}

		r.Log.Printf("Successfully created new PatientYearlyBenefitClaim with ID: %d for PatientID: %d, YearlyBenefitClaimID: %d", newPatientYearlyBenefitClaim.ID, patient.ID, yearlyBenefitClaim.ID)
		return &newPatientYearlyBenefitClaim, nil
	}

	r.Log.Printf("Database error checking PatientYearlyBenefitClaim for PatientID: %d, YearlyBenefitClaimID: %d: %v", patient.ID, yearlyBenefitClaim.ID, err)
	return nil, err
}

func (r *PatientYearlyBenefitClaimRepository) GetByPatientYearlyBenefitClaimID(db *gorm.DB, patientId, yearlyBenefitClaimId uint) (*entity.PatientYearlyBenefitClaim, error) {
	var patientYearlyBenefitClaim entity.PatientYearlyBenefitClaim

	baseQuery := db.Model(&entity.PatientYearlyBenefitClaim{})

	err := baseQuery.
		Where("patient_id = ? ", patientId).
		Where("yearly_benefit_claim_id = ? ", yearlyBenefitClaimId).
		Preload("Patient").
		Preload("YearlyBenefitClaim").
		Take(&patientYearlyBenefitClaim).Error
	if err != nil {
		return nil, err
	}

	return &patientYearlyBenefitClaim, nil
}

func (r *PatientYearlyBenefitClaimRepository) ChangeRemainingClaimPatientID(db *gorm.DB, patientId uint, yearlyClaimRemaining float64) error {
	baseQuery := db.Model(&entity.PatientYearlyBenefitClaim{})
	err := baseQuery.
		Where("patient_id = ? ", patientId).
		Update("yearly_claim_remaining", &yearlyClaimRemaining).
		Error
	if err != nil {
		return err
	}
	return nil
}
