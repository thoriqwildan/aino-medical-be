package helper

import (
	"math"
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

func ProRateRemainingMonthsFraction(t time.Time) float64 {
	rem := 13 - int(t.Month()) // exclude current month
	return float64(rem) / 12.0
}

func CalculateProrateYearlyMax(yearlyMax float64, prorate float64) float64 {
	prorateYearlyMax := yearlyMax * (prorate / 100)
	finalYearlyMax := yearlyMax - prorateYearlyMax
	return finalYearlyMax
}

func ProRateRemainingMonthsPercent(now, join time.Time) float64 {
	if join.IsZero() {
		return 0
	}
	if now.Location() != join.Location() {
		now = now.In(join.Location())
	}
	if now.Before(join) {
		return 0
	}
	if !now.Before(join.AddDate(1, 0, 0)) {
		return 100.0
	}
	return ProRateRemainingMonthsFraction(now)
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
		FindInBatches(&batch, 100, func(txBatch *gorm.DB, _ int) error {
			for _, emp := range batch {
				if emp == nil || emp.Patient.ID == 0 || emp.JoinDate.IsZero() {
					continue
				}
				if isAnniversary(emp.JoinDate, now) {
					if err := tx.Model(&entity.PatientBenefit{}).
						Where("patient_id = ?", emp.Patient.ID).
						Update("remaining_plafond", 0).
						Error; err != nil {
						return err
					}
				}
				for _, benefit := range emp.Patient.Benefits {
					if benefit.YearlyMax != nil {
						if err := tx.Model(&entity.PatientBenefit{}).
							Where("benefit_id = ?", benefit.ID).
							Where("patient_id = ?", emp.Patient.ID).
							Update("yearly_max", CalculateProrateYearlyMax(*benefit.YearlyMax, emp.ProRate)).Error; err != nil {
							return err
						}

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

func round2(x float64) float64 { return math.Round(x*100) / 100 }
