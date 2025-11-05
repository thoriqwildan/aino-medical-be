package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
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

func (r *PatientBenefitRepository) GetAll(db *gorm.DB, request *model.SearchPagingQuery) ([]*entity.PatientBenefit, int64, error) {
	var (
		rows  []*entity.PatientBenefit
		total int64
	)

	base := db.Model(&entity.PatientBenefit{}).
		Joins("JOIN patients p ON p.id = patient_benefits.patient_id")

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
		Distinct("patient_benefits.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := base.
		Scopes(filterScope(request)).
		Select("patient_benefits.*").
		Preload("Patient").
		Preload("Benefit").
		Preload("Claims").
		Order("p.name ASC")

	if request != nil {
		if request.Limit > 0 {
			q = q.Limit(request.Limit)
		}
		if request.Page > 0 && request.Limit > 0 {
			q = q.Offset((request.Page - 1) * request.Limit)
		}
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *PatientBenefitRepository) FindOrCreatePatient(db *gorm.DB, patient *entity.Patient) (*entity.Patient, error) {
	var result entity.Patient
	errTake := db.Model(entity.Patient{}).Where("name = ?", patient.Name).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		r.Log.Error("Error when find patient in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil
	}

	errCreate := db.Create(patient).Error
	if errCreate != nil {
		r.Log.Error("Error when create patient in method find or create")
		return nil, errCreate
	}

	return patient, nil
}

func (r *PatientBenefitRepository) GetByPatientBenefitID(db *gorm.DB, patientId, benefitId uint) (*entity.PatientBenefit, error) {
	var patientBenefit entity.PatientBenefit

	baseQuery := db.Model(&entity.PatientBenefit{})

	err := baseQuery.
		Where("patient_id = ? ", patientId).
		Where("benefit_id = ? ", benefitId).
		Preload("Patient").
		Preload("Benefit").
		Preload("Claims").
		Take(&patientBenefit).Error
	if err != nil {
		return nil, err
	}

	return &patientBenefit, nil
}

func (r *PatientBenefitRepository) ChangeRemainingPlafondByPatientID(db *gorm.DB, patientId uint, remainingPlafond float64) error {
	baseQuery := db.Model(&entity.PatientBenefit{})
	err := baseQuery.
		Where("patient_id = ? ", patientId).
		Update("remaining_plafond", &remainingPlafond).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PatientBenefitRepository) FindOrCreate(
	db *gorm.DB,
	patient *entity.Patient,
	benefit *entity.Benefit,
	initialPlafond *float64,
	startDate time.Time,
	prorate float64,
) (*entity.PatientBenefit, error) {
	// TODO:
	if benefit.LimitationType == entity.LimitationTypePerDay || benefit.LimitationType == entity.LimitationTypePerIncident {
		patientBenefit := entity.PatientBenefit{
			PatientID:        patient.ID,
			BenefitID:        benefit.ID,
			InitialPlafond:   initialPlafond,
			RemainingPlafond: nil,
			StartDate:        startDate,
		}
		errCreate := db.Model(entity.PatientBenefit{}).Create(&patientBenefit).Error
		if errCreate != nil {
			return nil, errCreate
		}
		return &patientBenefit, nil
	}
	
	var patientBenefit entity.PatientBenefit

	err := db.Where("patient_id = ? AND benefit_id = ?", patient.ID, benefit.ID).First(&patientBenefit).Error

	if err == nil {
		r.Log.Printf("PatientBenefit found for PatientID: %d, BenefitID: %d", patient.ID, benefit.ID)
		return &patientBenefit, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.Log.Printf("PatientBenefit not found for PatientID: %d, BenefitID: %d. Creating new record...", patient.ID, benefit.ID)
		var remainingPlafond *float64
		if initialPlafond != nil {
			var results float64
			if patient.Employee != nil {
				results = helper.CalculateProratePlafond(*initialPlafond, patient.Employee.ProRate)
			} else {
				results = helper.CalculateProratePlafond(*initialPlafond, prorate)
			}
			remainingPlafond = &results
		}
		newPatientBenefit := entity.PatientBenefit{
			PatientID:        patient.ID,
			BenefitID:        benefit.ID,
			InitialPlafond:   initialPlafond,
			RemainingPlafond: remainingPlafond,
			StartDate:        startDate,
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
		return errors.New("error cannot reduction remaining plafond because of nil plafond value")
	}
	return db.Save(patientBenefit).Error
}
