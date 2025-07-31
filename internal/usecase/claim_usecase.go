package usecase

import (
	"context"
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

type ClaimUseCase struct {
	Repository *repository.ClaimRepository
	PatientBenefitRepository *repository.PatientBenefitRepository
	Log *logrus.Logger
	DB *gorm.DB
	Validate *validator.Validate
}

func NewClaimUseCase(repo *repository.ClaimRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger, patientBenefitRepository *repository.PatientBenefitRepository) *ClaimUseCase {
	return &ClaimUseCase{
		Repository: repo,
		DB: db,
		Validate: validate,
		Log: log,
		PatientBenefitRepository: patientBenefitRepository,
	}
}

func (uc *ClaimUseCase) Create(ctx context.Context, request *model.ClaimRequest) (*model.ClaimResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation error in ClaimRequest")
		return nil, err
	}

	now := time.Now()
	SLA := helper.DetermineSLAStatus(now)
	startDateOfCurrentYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())

	benefit := &entity.Benefit{}
	if err := uc.Repository.GetBenefitByCode(tx, benefit, request.BenefitCode); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.Error("Benefit not found")
			return nil, err
		}
		uc.Log.WithError(err).Error("Failed to get benefit by code")
		return nil, err
	}

	patient := &entity.Patient{}
	if err := uc.Repository.GetPatientByID(tx, patient, request.PatientID); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.Error("Patient not found")
			return nil, fiber.NewError(fiber.StatusNotFound, "Patient not found")
		}
		uc.Log.WithError(err).Error("Failed to find patient")
		return nil, err
	}

	if patient.PlanTypeID != benefit.PlanTypeID {
		uc.Log.Error("Patient's plan type does not match benefit's plan type")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Patient's plan type does not match benefit's plan type")
	}

	patientBenefit, err := uc.PatientBenefitRepository.FindOrCreate(tx, patient.ID, benefit.ID, float64(benefit.Plafond), startDateOfCurrentYear)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find or create patient benefit")
		return nil, err
	}

	claim := &entity.Claim{
		PatientID: 	 request.PatientID,
		PatientBenefitID: patientBenefit.ID,
		ClaimAmount: request.ClaimAmount,
		SLA: &SLA,
		TransactionStatus: entity.TransactionStatusPending,
	}

	if patientBenefit.RemainingPlafond < request.ClaimAmount {
		claim.ClaimStatus = entity.ClaimStatusOverPlafond
		claim.ApprovedAmount = &patientBenefit.RemainingPlafond
	} else {
		claim.ClaimStatus = entity.ClaimStatusOnPlafond
		claim.ApprovedAmount = &request.ClaimAmount
	}

	if patient.FamilyMemberID != nil {
		claim.EmployeeID = patient.FamilyMember.EmployeeID
	} else {
		claim.EmployeeID = *patient.EmployeeID
	}

	if err := uc.PatientBenefitRepository.BalanceReduction(tx, patientBenefit, claim.ClaimAmount); err != nil {
		uc.Log.WithError(err).Error("Failed to reduce patient benefit balance")
		if err == gorm.ErrInvalidData {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient benefit balance")
		}
		return nil, err
	}

	if err := uc.Repository.Create(tx, claim); err != nil {
		uc.Log.WithError(err).Error("Failed to create claim")
		return nil, err
	}	

	if err := uc.Repository.GetByID(tx, claim, claim.ID); err != nil {
		uc.Log.WithError(err).Error("Failed to retrieve claim by ID after creation")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction in CreateClaim")
		return nil, err
	}

	return converter.ClaimToResponse(claim), nil
}

func (uc *ClaimUseCase) GetPatient(ctx context.Context, request *model.PagingQuery) ([]model.PatientResponse, int64, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	patients, total, err := uc.Repository.GetPatients(tx, request)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to get patients")
		return nil, 0, err
	}

	if len(patients) == 0 {
		uc.Log.Error("No patients found")
		return nil, 0, fiber.NewError(fiber.StatusNotFound, "No patients found")
	}

	responses := make([]model.PatientResponse, len(patients))
	for i, p := range patients {
		responses[i] = *converter.PatientToResponse(&p)
	}
	return responses, total, nil
}

func (uc *ClaimUseCase) GetBenefit(ctx context.Context, request *model.PagingQuery, patientId uint) ([]model.BenefitResponse, int64, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	patient := &entity.Patient{}
	if err := uc.Repository.GetPatientByID(tx, patient, patientId); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("patientId", patientId).Error("Patient not found in GetBenefit")
			return nil, 0, fiber.NewError(fiber.StatusNotFound, "Patient not found")
		}
		uc.Log.WithError(err).Error("Failed to get patient by ID in GetBenefit")
		return nil, 0, err
	}

	benefits, total, err := uc.Repository.GetBenefits(tx, request, patient.PlanType.ID)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to get benefits")
		return nil, 0, err
	}

	responses := make([]model.BenefitResponse, len(benefits))
	for i, b := range benefits {
		responses[i] = *converter.BenefitToResponse(&b)
	}
	return responses, total, nil
}