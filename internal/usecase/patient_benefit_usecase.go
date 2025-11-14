package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type PatientBenefitUsecase struct {
	Repository *repository.PatientBenefitRepository
	Log        *logrus.Logger
	DB         *gorm.DB
	Validate   *validator.Validate
}

func NewPatientBenefitUsecase(repository *repository.PatientBenefitRepository, log *logrus.Logger, DB *gorm.DB, validate *validator.Validate) *PatientBenefitUsecase {
	return &PatientBenefitUsecase{Repository: repository, Log: log, DB: DB, Validate: validate}
}

func startOfYearUTC(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func (p PatientBenefitUsecase) checkPatientBenefit(tx *gorm.DB, params *model.PatientBenefitParams) (*entity.Patient, *entity.Benefit, error) {
	var benefit entity.Benefit
	if err := tx.Model(&entity.Benefit{}).First(&benefit, "id = ?", params.BenefitID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fiber.NewError(fiber.StatusNotFound, "benefit not found")
		}
		return nil, nil, fiber.NewError(fiber.StatusInternalServerError, "error when find benefit")
	}

	var patient entity.Patient
	if err := tx.Model(&entity.Patient{}).Preload("Employee").First(&patient, "id = ?", params.PatientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &benefit, fiber.NewError(fiber.StatusNotFound, "patient not found")
		}
		return nil, &benefit, fiber.NewError(fiber.StatusInternalServerError, "error when find patient")
	}

	return &patient, &benefit, nil
}

func (p *PatientBenefitUsecase) Create(ctx context.Context, request *model.PatientBenefitParams) (*model.PatientBenefitResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate create patient benefit: %v", err)
		return nil, err
	}

	patient, benefit, err := p.checkPatientBenefit(tx, request)
	if err != nil {
		return nil, err
	}

	startDate := startOfYearUTC(time.Now())
	patientBenefit, err := p.Repository.FindOrCreate(tx, patient, benefit, benefit.Plafond, startDate, patient.Employee.ProRate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit create patient benefit")
		return nil, err
	}
	return converter.PatientBenefitToResponse(patientBenefit), nil
}

func (p *PatientBenefitUsecase) GetByPatientBenefitID(ctx context.Context, request *model.PatientBenefitParams) (*model.PatientBenefitResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate find patient benefit: %v", err)
		return nil, err
	}

	patient, benefit, err := p.checkPatientBenefit(tx, request)
	if err != nil {
		return nil, err
	}

	pb, err := p.Repository.GetByPatientBenefitID(tx, patient.ID, benefit.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("error cause patient benefit not found: %v", err))
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit find patient benefit")
		return nil, err
	}
	return converter.PatientBenefitToResponse(pb), nil
}

func (p PatientBenefitUsecase) GetAll(ctx context.Context, request *model.SearchPagingQuery) ([]*model.PatientBenefitResponse, int64, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate find all patient benefit: %v", err)
		return nil, 0, err
	}

	pbs, count, err := p.Repository.GetAll(tx, request)
	if err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").WithError(err).Error("find all patient benefit")
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit find all patient benefit")
		return nil, 0, err
	}

	res := make([]*model.PatientBenefitResponse, len(pbs))
	for i, pb := range pbs {
		res[i] = converter.PatientBenefitToResponse(pb)
	}
	return res, count, nil
}

func (p *PatientBenefitUsecase) Update(ctx context.Context, params *model.PatientBenefitParams, request *model.UpdatePatientBenefitRequest) (*model.PatientBenefitResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate params update patient benefit: %v", err)
		return nil, err
	}
	if err := p.Validate.Struct(request); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate body update patient benefit: %v", err)
		return nil, err
	}

	patient, benefit, err := p.checkPatientBenefit(tx, params)
	if err != nil {
		return nil, err
	}

	startDate := startOfYearUTC(time.Now())
	pb, err := p.Repository.FindOrCreate(tx, patient, benefit, benefit.Plafond, startDate, patient.Employee.ProRate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if request.InitialPlafond != nil {
		pb.InitialPlafond = request.InitialPlafond
	}
	if request.EndDate != nil {
		pb.EndDate = (*time.Time)(request.EndDate)
	}
	if request.StartDate != nil {
		t := *request.StartDate
		pb.StartDate = (time.Time)(t)
	}

	if pb.InitialPlafond != nil {
		results := helper.CalculateProratePlafond(*pb.InitialPlafond, patient.Employee.ProRate)
		pb.RemainingPlafond = &results
	}

	switch request.Status {
	case "active":
		pb.Status = entity.PatientBenefitStatusActive
	case "exhausted":
		pb.Status = entity.PatientBenefitStatusExhausted
	case "expired":
		pb.Status = entity.PatientBenefitStatusExpired
	}

	if err := tx.Save(pb).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit update patient benefit")
		return nil, err
	}
	return converter.PatientBenefitToResponse(pb), nil
}

func (p *PatientBenefitUsecase) ResetRemainingPlafondByPatientBenefitID(ctx context.Context, params *model.PatientBenefitParams) (*model.PatientBenefitResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate params reset remaining plafond: %v", err)
		return nil, err
	}

	patient, benefit, err := p.checkPatientBenefit(tx, params)
	if err != nil {
		return nil, err
	}

	startDate := startOfYearUTC(time.Now())
	pb, err := p.Repository.FindOrCreate(tx, patient, benefit, benefit.Plafond, startDate, patient.Employee.ProRate)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	pb.RemainingPlafond = pb.InitialPlafond
	if err := tx.Save(pb).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit reset remaining plafond by id")
		return nil, err
	}
	return converter.PatientBenefitToResponse(pb), nil
}

func (p *PatientBenefitUsecase) ResetRemainingPlafondByPatientID(ctx context.Context, patientId uint) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var list []*entity.PatientBenefit
		return tx.Model(&entity.PatientBenefit{}).
			Where("patient_id = ?", patientId).
			FindInBatches(&list, 100, func(tx *gorm.DB, batch int) error {
				for _, pb := range list {
					pb.RemainingPlafond = pb.InitialPlafond
					if err := tx.Save(pb).Error; err != nil {
						return err
					}
				}
				return nil
			}).Error
	})
}

func (p *PatientBenefitUsecase) ResetAllRemainingPlafond(ctx context.Context) error {
	return p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var list []*entity.PatientBenefit
		return tx.Model(&entity.PatientBenefit{}).
			FindInBatches(&list, 100, func(tx *gorm.DB, batch int) error {
				for _, pb := range list {
					pb.RemainingPlafond = pb.InitialPlafond
					if err := tx.Save(pb).Error; err != nil {
						return err
					}
				}
				return nil
			}).Error
	})
}

func (p *PatientBenefitUsecase) Delete(ctx context.Context, params *model.PatientBenefitParams) error {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := p.Validate.Struct(params); err != nil {
		p.Log.WithField("usecase", "PatientBenefitUsecase").Errorf("validate params delete patient benefit: %v", err)
		return err
	}

	patient, benefit, err := p.checkPatientBenefit(tx, params)
	if err != nil {
		return err
	}

	pb, err := p.Repository.GetByPatientBenefitID(tx, patient.ID, benefit.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("error cause patient benefit not found: %v", err))
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if pb == nil || pb.ID == 0 {
		return fiber.NewError(fiber.StatusNotFound, "patient benefit not found")
	}

	if err := tx.Delete(pb).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when delete patient benefit: %s", err.Error()))
	}

	if err := tx.Commit().Error; err != nil {
		p.Log.WithError(err).Error("commit delete patient benefit")
		return err
	}
	return nil
}
