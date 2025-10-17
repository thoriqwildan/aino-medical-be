package helper

import (
	"errors"
	"math"
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

func CalculateProratePlafond(plafond float64, prorate float64) float64 {
	return plafond * (prorate / 100)
}

func monthsRemainingFromMonth(m time.Month, includeCurrent bool) int {
	if includeCurrent {
		return 13 - int(m) // include current month: Sep (9) -> 4 (Sep..Dec)
	}
	return 12 - int(m) // exclude current month: Sep (9) -> 3 (Oct..Dec)
}

func ProRateRemainingMonthsPercent(now, join time.Time) float64 {
	if join.IsZero() {
		return 0
	}

	if now.Location() != join.Location() {
		now = now.In(join.Location())
	}

	if !now.Before(join.AddDate(1, 0, 0)) {
		return 100.0
	}

	rem := monthsRemainingFromMonth(join.Month(), true)
	if rem < 0 {
		rem = 0
	}
	percent := (float64(rem) / 12.0) * 100.0
	return math.Round(percent*100) / 100
}

func ResetBenefitProRateDaily(db *gorm.DB) error {
	now := time.Now()
	var employees []*entity.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}

	for _, employee := range employees {
		join := employee.JoinDate
		if err := db.Model(&entity.Employee{}).
			Where("id = ?", employee.ID).
			Update("pro_rate", ProRateRemainingMonthsPercent(now, join)).Error; err != nil {
			return err
		}
	}
	return nil
}

func ResetPatientBenefitRemainingPlafondDaily(db *gorm.DB) error {
	now := time.Now()

	tx := db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	batch := make([]*entity.Employee, 100)
	if err := db.Model(&entity.Employee{}).
		Preload("Patient").
		Preload("Patient.Benefits").
		Preload("Patient.Benefits.YearlyBenefitClaim").
		FindInBatches(&batch, 100, func(txBatch *gorm.DB, _ int) error {
			for _, emp := range batch {
				for _, benefit := range emp.Patient.Benefits {
					if emp == nil || emp.Patient.ID == 0 || emp.JoinDate.IsZero() {
						continue
					}
					if isAnniversary(emp.JoinDate, now) {
						if benefit.YearlyBenefitClaimID != nil {
							var patientYearlyBenefitClaim entity.PatientYearlyBenefitClaim
							err := tx.Model(&entity.PatientYearlyBenefitClaim{}).
								Where("patient_id = ?", emp.Patient.ID).
								Where("yearly_benefit_claim_id = ?", benefit.YearlyBenefitClaimID).
								Take(&patientYearlyBenefitClaim).Error
							if err != nil {
								if errors.Is(err, gorm.ErrRecordNotFound) {
									continue
								}
								return err
							}
							if benefit.YearlyBenefitClaim != nil {
								patientYearlyBenefitClaim.YearlyClaimRemaining = benefit.YearlyBenefitClaim.YearlyClaim
							}
							tx.Save(&patientYearlyBenefitClaim)
						}
						var patientBenefit entity.PatientBenefit
						err := tx.Model(entity.PatientBenefit{}).
							Where("patient_id = ?", emp.Patient.ID).
							Where("benefit_id = ?", benefit.ID).
							Take(&patientBenefit).Error
						if err != nil {
							if errors.Is(err, gorm.ErrRecordNotFound) {
								continue
							}
							return err
						}
						if patientBenefit.InitialPlafond != nil {
							patientBenefit.RemainingPlafond = patientBenefit.InitialPlafond
						} else {
							patientBenefit.RemainingPlafond = nil
						}
						db.Save(&patientBenefit)
					}
				}
			}
			return nil
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
