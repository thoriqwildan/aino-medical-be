package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/model/converter"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type PatientYearlyBenefitClaimUsecase struct {
	Repository *repository.PatientYearlyBenefitClaimRepository
	Log        *logrus.Logger
	DB         *gorm.DB
	Validate   *validator.Validate
}

func NewPatientYearlyBenefitClaimUsecase(repo *repository.PatientYearlyBenefitClaimRepository, log *logrus.Logger, DB *gorm.DB, validate *validator.Validate) *PatientYearlyBenefitClaimUsecase {
	return &PatientYearlyBenefitClaimUsecase{Repository: repo, Log: log, DB: DB, Validate: validate}
}

func (p PatientYearlyBenefitClaimUsecase) checkPatientYearlyBenefitClaim(tx *gorm.DB, params *model.PatientYearlyBenefitClaimParams) (*entity.Patient, *entity.YearlyBenefitClaim, error) {
	var yearlyBenefitClaim entity.YearlyBenefitClaim
	if err := tx.Model(&entity.YearlyBenefitClaim{}).First(&yearlyBenefitClaim, "id = ?", params.YearlyBenefitClaimID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fiber.NewError(fiber.StatusNotFound, "yearly benefit claim not found")
		}
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "error when find yearly benefit claim")
	}

	var patient entity.Patient
	if err := tx.Model(&entity.Patient{}).Preload("Employee").First(&patient, "id = ?", params.PatientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &yearlyBenefitClaim, fiber.NewError(fiber.StatusNotFound, "patient not found")
		}
		return nil, &yearlyBenefitClaim, fiber.NewError(fiber.StatusInternalServerError, "error when find patient")
	}

	return &patient, &yearlyBenefitClaim, nil
}

func (p *PatientYearlyBenefitClaimUsecase) Create(ctx context.Context, request *model.PatientYearlyBenefitClaimParams) (*model.PatientYearlyBenefitClaimResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate create patient yearly benefit claim: %v", err)
		return nil, err
	}

	patient, yearlyBenefitClaim, err := p.checkPatientYearlyBenefitClaim(tx, request)
	if err != nil {
		return nil, err
	}

	patientYearlyBenefitClaim, err := p.Repository.FindOrCreate(tx, *patient, *yearlyBenefitClaim)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit create patient yearly benefit claim")
		return nil, err
	}
	return converter.PatientYearlyBenefitClaimToResponse(patientYearlyBenefitClaim), nil
}

func (p *PatientYearlyBenefitClaimUsecase) GetByPatientYearlyBenefitClaimID(ctx context.Context, request *model.PatientYearlyBenefitClaimParams) (*model.PatientYearlyBenefitClaimResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate find patient yearly benefit claim: %v", err)
		return nil, err
	}

	patient, yearlyBenefitClaim, err := p.checkPatientYearlyBenefitClaim(tx, request)
	if err != nil {
		return nil, err
	}

	patientYearlyBenefitClaim, err := p.Repository.GetByPatientYearlyBenefitClaimID(tx, patient.ID, yearlyBenefitClaim.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("error cause patient yearly benefit claim not found: %v", err))
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit find patient yearly benefit claim")
		return nil, err
	}
	return converter.PatientYearlyBenefitClaimToResponse(patientYearlyBenefitClaim), nil
}

func (p PatientYearlyBenefitClaimUsecase) GetAll(ctx context.Context, request *model.SearchPagingQuery) ([]*model.PatientYearlyBenefitClaimResponse, int64, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate find all patient yearly benefit claim: %v", err)
		return nil, 0, err
	}

	patientYearlyBenefitClaims, count, err := p.Repository.GetAll(tx, request)
	if err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").WithError(err).Error("find all patient yearly benefit claim")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit find all patient yearly benefit claim")
		return nil, 0, err
	}

	res := make([]*model.PatientYearlyBenefitClaimResponse, len(patientYearlyBenefitClaims))
	for i, pb := range patientYearlyBenefitClaims {
		res[i] = converter.PatientYearlyBenefitClaimToResponse(pb)
	}
	return res, count, nil
}

func (p *PatientYearlyBenefitClaimUsecase) Update(ctx context.Context, params *model.PatientYearlyBenefitClaimParams) (*model.PatientYearlyBenefitClaimResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate params update patient yearly benefit claim: %v", err)
		return nil, err
	}

	patient, yearlyBenefitClaim, err := p.checkPatientYearlyBenefitClaim(tx, params)
	if err != nil {
		return nil, err
	}
	patientYearlyBenefitClaim, err := p.Repository.FindOrCreate(tx, *patient, *yearlyBenefitClaim)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	patientYearlyBenefitClaim.YearlyClaimRemaining = helper.CalculateProratePlafond(yearlyBenefitClaim.YearlyClaim, patient.Employee.ProRate)
	if err := tx.Save(patientYearlyBenefitClaim).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit update patient yearly benefit claim")
		return nil, err
	}
	return converter.PatientYearlyBenefitClaimToResponse(patientYearlyBenefitClaim), nil
}

func (p *PatientYearlyBenefitClaimUsecase) ResetRemainingPlafondByPatientYearlyBenefitClaimID(ctx context.Context, params *model.PatientYearlyBenefitClaimParams) (*model.PatientYearlyBenefitClaimResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate params reset remaining plafond: %v", err)
		return nil, err
	}

	patient, yearlyBenefitClaim, err := p.checkPatientYearlyBenefitClaim(tx, params)
	if err != nil {
		return nil, err
	}

	patientYearlyBenefitClaim, err := p.Repository.FindOrCreate(tx, *patient, *yearlyBenefitClaim)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	patientYearlyBenefitClaim.YearlyClaimRemaining = yearlyBenefitClaim.YearlyClaim
	if err := tx.Save(patientYearlyBenefitClaim).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit reset remaining plafond by id")
		return nil, err
	}
	return converter.PatientYearlyBenefitClaimToResponse(patientYearlyBenefitClaim), nil
}

func (p *PatientYearlyBenefitClaimUsecase) ResetRemainingPlafondByPatientID(ctx context.Context, patientId uint) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var list []*entity.PatientYearlyBenefitClaim
		return tx.Model(&entity.PatientYearlyBenefitClaim{}).
			Preload("YearlyBenefitClaim").
			Where("patient_id = ?", patientId).
			FindInBatches(&list, 100, func(tx *gorm.DB, batch int) error {
				for _, patientYearlyBenefitClaim := range list {
					patientYearlyBenefitClaim.YearlyClaimRemaining = patientYearlyBenefitClaim.YearlyBenefitClaim.YearlyClaim
					if err := tx.Save(patientYearlyBenefitClaim).Error; err != nil {
						return err
					}
				}
				return nil
			}).Error
	})
}

func (p *PatientYearlyBenefitClaimUsecase) ResetAllRemainingPlafond(ctx context.Context) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var list []*entity.PatientYearlyBenefitClaim
		return tx.Model(&entity.PatientYearlyBenefitClaim{}).
			Preload("YearlyBenefitClaim").
			FindInBatches(&list, 100, func(tx *gorm.DB, batch int) error {
				for _, patientYearlyBenefitClaim := range list {
					patientYearlyBenefitClaim.YearlyClaimRemaining = patientYearlyBenefitClaim.YearlyBenefitClaim.YearlyClaim
					if err := tx.Save(patientYearlyBenefitClaim).Error; err != nil {
						return err
					}
				}
				return nil
			}).Error
	})
}

func (p *PatientYearlyBenefitClaimUsecase) Delete(ctx context.Context, params *model.PatientYearlyBenefitClaimParams) error {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientYearlyBenefitClaim").Errorf("validate params delete patient yearly benefit claim: %v", err)
		return err
	}

	patient, yearlyBenefitClaim, err := p.checkPatientYearlyBenefitClaim(tx, params)
	if err != nil {
		return err
	}

	patientYearlyBenefitClaim, err := p.Repository.GetByPatientYearlyBenefitClaimID(tx, patient.ID, yearlyBenefitClaim.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("error cause patient yearly benefit claim not found: %v", err))
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if patientYearlyBenefitClaim == nil || patientYearlyBenefitClaim.ID == 0 {
		return fiber.NewError(fiber.StatusNotFound, "patient yearly benefit claim not found")
	}

	if err := tx.Delete(patientYearlyBenefitClaim).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when delete patient yearly benefit claim: %s", err.Error()))
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit delete patient yearly benefit claim")
		return err
	}
	return nil
}
