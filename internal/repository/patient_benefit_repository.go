package repository

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
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
	patient *entity.Patient,
	benefit *entity.Benefit,
	initialPlafond *float64,
	startDate time.Time,
	prorate float64,
) (*entity.PatientBenefit, error) {
	var patientBenefit entity.PatientBenefit

	err := db.Where("patient_id = ? AND benefit_id = ?", patient.ID, benefit.ID).First(&patientBenefit).Error

	if err == nil {
		r.Log.Printf("PatientBenefit found for PatientID: %d, BenefitID: %d", patient.ID, benefit.ID)
		return &patientBenefit, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.Log.Printf("PatientBenefit not found for PatientID: %d, BenefitID: %d. Creating new record...", patient.ID, benefit.ID)

		newPatientBenefit := entity.PatientBenefit{
			PatientID:        patient.ID,
			BenefitID:        benefit.ID,
			InitialPlafond:   initialPlafond,
			RemainingPlafond: initialPlafond,
			StartDate:        startDate,
		}

		if benefit.YearlyMax != nil {
			yearlyMax := helper.CalculateProrateYearlyMax(*benefit.YearlyMax, prorate)
			newPatientBenefit.YearlyMax = &yearlyMax
		}

		createErr := db.Create(&newPatientBenefit).Error
		if createErr != nil {
			r.Log.Printf("Error creating PatientBenefit for PatientID: %d, BenefitID: %d: %v", patient.ID, benefit.ID, createErr)
			return nil, createErr
		}

		r.Log.Printf("Successfully created new PatientBenefit with ID: %d for PatientID: %d, BenefitID: %d", newPatientBenefit.ID, patient.ID, benefit.ID)
		return &newPatientBenefit, nil
	}

	r.Log.Printf("Database error checking PatientBenefit for PatientID: %d, BenefitID: %d: %v", patient.ID, benefit.ID, err)
	return nil, err
}

func (r *PatientBenefitRepository) BalanceReduction(db *gorm.DB, patientBenefit *entity.PatientBenefit, amount float64) error {

	if patientBenefit.RemainingPlafond != nil {
		*patientBenefit.RemainingPlafond -= amount
		if *patientBenefit.RemainingPlafond < 0 {
			return gorm.ErrInvalidData
		}
	} else {
		return errors.New("error cannot reductio remaining plafond because of nil plafond value")
	}
	return db.Save(patientBenefit).Error
}
