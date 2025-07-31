package repository

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

type PatientBenefitRepository struct {
	Repository[entity.PatientBenefit]
	Log *logrus.Logger
}

func NewPatientBenefitRepository(log *logrus.Logger) *PatientBenefitRepository {
	return &PatientBenefitRepository{
		Log: log,
	}
}

func (r *PatientBenefitRepository) FindOrCreate(
	db *gorm.DB,
	patientID uint,
	benefitID uint,
	initialPlafond float64,
	startDate time.Time,
) (*entity.PatientBenefit, error) {
	var patientBenefit entity.PatientBenefit

	err := db.Where("patient_id = ? AND benefit_id = ?", patientID, benefitID).First(&patientBenefit).Error

	if err == nil {
		r.Log.Printf("PatientBenefit found for PatientID: %d, BenefitID: %d", patientID, benefitID)
		return &patientBenefit, nil
	}

	if err == gorm.ErrRecordNotFound {
		r.Log.Printf("PatientBenefit not found for PatientID: %d, BenefitID: %d. Creating new record...", patientID, benefitID)

		newPatientBenefit := entity.PatientBenefit{
			PatientID:        patientID,
			BenefitID:        benefitID,
			InitialPlafond:   initialPlafond,
			RemainingPlafond: initialPlafond,
			StartDate:        startDate,
		}

		createErr := db.Create(&newPatientBenefit).Error
		if createErr != nil {
			r.Log.Printf("Error creating PatientBenefit for PatientID: %d, BenefitID: %d: %v", patientID, benefitID, createErr)
			return nil, createErr
		}

		r.Log.Printf("Successfully created new PatientBenefit with ID: %d for PatientID: %d, BenefitID: %d", newPatientBenefit.ID, patientID, benefitID)
		return &newPatientBenefit, nil
	}

	r.Log.Printf("Database error checking PatientBenefit for PatientID: %d, BenefitID: %d: %v", patientID, benefitID, err)
	return nil, err
}

func (r *PatientBenefitRepository) BalanceReduction(db *gorm.DB, patientBenefit *entity.PatientBenefit, amount float64) error {
	patientBenefit.RemainingPlafond -= amount
	if patientBenefit.RemainingPlafond < 0 {
		return gorm.ErrInvalidData
	}

	return db.Save(patientBenefit).Error
}