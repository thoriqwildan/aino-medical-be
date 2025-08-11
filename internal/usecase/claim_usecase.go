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
	BenefitRepository *repository.BenefitRepository
	Log *logrus.Logger
	DB *gorm.DB
	Validate *validator.Validate
}

func NewClaimUseCase(repo *repository.ClaimRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger, patientBenefitRepository *repository.PatientBenefitRepository, benefitRepository *repository.BenefitRepository) *ClaimUseCase {
	return &ClaimUseCase{
		Repository: repo,
		DB: db,
		Validate: validate,
		Log: log,
		PatientBenefitRepository: patientBenefitRepository,
		BenefitRepository: benefitRepository,
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
    // ... error handling ...
  }
  
  // PANGGIL METHOD REPOSITORY YANG BARU
  benefits, remainingPlafondMap, total, err := uc.Repository.GetBenefitsWithPlafond(tx, request, patient.PlanTypeID, patientId)
  if err != nil {
    uc.Log.WithError(err).Error("Failed to get benefits with plafond")
    return nil, 0, err
  }

  // Lakukan konversi dengan data tambahan
  responses := make([]model.BenefitResponse, len(benefits))
  for i, b := range benefits {
    response := converter.BenefitToResponse(&b)
    
    // Cek apakah ada remaining plafond untuk benefit ini
    if rp, ok := remainingPlafondMap[b.ID]; ok {
      response.RemainingPlafond = &rp
    }
    
    responses[i] = *response
  }
  
  // Commit transaksi
  if err := tx.Commit().Error; err != nil {
      uc.Log.WithError(err).Error("Failed to commit transaction in GetBenefit")
      return nil, 0, err
  }
  
  return responses, total, nil
}

func (uc *ClaimUseCase) UpdateClaim(ctx context.Context, request *model.UpdateClaimRequest) (*model.ClaimResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation error in UpdateClaimRequest")
		return nil, err
	}

	now := time.Now()
	SLA := helper.DetermineSLAStatus(now)
	startDateOfCurrentYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())

	claim := &entity.Claim{}
	if err := uc.Repository.GetByID(tx, claim, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("id", request.ID).Error("Claim not found in UpdateClaim")
			return nil, fiber.NewError(fiber.StatusNotFound, "Claim not found")
		}
		uc.Log.WithError(err).Error("Failed to get claim by ID in UpdateClaim")
		return nil, err
	}

	benefit := &entity.Benefit{}
	if err := uc.BenefitRepository.GetById(tx, claim.PatientBenefit.BenefitID, benefit); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("benefitId", claim.PatientBenefit.BenefitID).Error("Benefit not found in UpdateClaim")
			return nil, fiber.NewError(fiber.StatusNotFound, "Benefit not found")
		}
		uc.Log.WithError(err).Error("Failed to get benefit by ID in UpdateClaim")
		return nil, err
	}

	patientBenefit, err := uc.PatientBenefitRepository.FindOrCreate(tx, claim.PatientID, benefit.ID, float64(benefit.Plafond), startDateOfCurrentYear)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find or create patient benefit in UpdateClaim")
		return nil, err
	}

	patientBenefit.RemainingPlafond += *claim.ApprovedAmount
	if err := uc.PatientBenefitRepository.BalanceReduction(tx, patientBenefit, request.ClaimAmount); err != nil {
		uc.Log.WithError(err).Error("Failed to reduce patient benefit balance in UpdateClaim")
		if err == gorm.ErrInvalidData {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient benefit balance")
		}
		return nil, err
	}

	claim.ClaimAmount = request.ClaimAmount
	claim.SLA = &SLA
	claim.TransactionTypeID = request.TransactionTypeID
	claim.TransactionStatus = entity.TransactionStatus(request.TransactionStatus)
	if patientBenefit.RemainingPlafond < request.ClaimAmount {
		claim.ClaimStatus = entity.ClaimStatusOverPlafond
		claim.ApprovedAmount = &patientBenefit.RemainingPlafond
	} else {
		claim.ClaimStatus = entity.ClaimStatusOnPlafond
		claim.ApprovedAmount = &request.ClaimAmount
	}
	claim.SubmissionDate = (*time.Time)(request.SubmissionDate)
	claim.City = request.City
	claim.Diagnosis = request.Diagnosis
	claim.MedicalFacilityName = request.MedicalFacility
	claim.DocLink = request.DocLink
	claim.TransactionDate = (*time.Time)(request.TransactionDate)

	if err := uc.Repository.Update(tx, claim); err != nil {
		uc.Log.WithError(err).Error("Failed to update claim")
		return nil, err
	}

	if err := uc.Repository.GetByID(tx, claim, claim.ID); err != nil {
		uc.Log.WithError(err).Error("Failed to retrieve claim by ID after update")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction in UpdateClaim")
		return nil, err
	}

	return converter.ClaimToResponse(claim), nil
}

func (uc *ClaimUseCase) GetClaim(ctx context.Context, id uint) (*model.ClaimResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claim := &entity.Claim{}
	if err := uc.Repository.GetByID(tx, claim, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("id", id).Error("Claim not found in GetClaim")
			return nil, fiber.NewError(fiber.StatusNotFound, "Claim not found")
		}
		uc.Log.WithError(err).Error("Failed to get claim by ID in GetClaim")
		return nil, err
	}

	response := converter.ClaimToResponse(claim)
	return response, nil
}

func (uc *ClaimUseCase) DeleteClaim(ctx context.Context, id uint) error {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claim := &entity.Claim{}
	if err := uc.Repository.GetByID(tx, claim, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("id", id).Error("Claim not found in DeleteClaim")
			return fiber.NewError(fiber.StatusNotFound, "Claim not found")
		}
		uc.Log.WithError(err).Error("Failed to get claim by ID in DeleteClaim")
		return err
	}

	patientBenefit := &entity.PatientBenefit{}
	if err := uc.PatientBenefitRepository.FindById(tx, patientBenefit, claim.PatientBenefitID); err != nil {
		if err == gorm.ErrRecordNotFound {
			uc.Log.WithField("patientBenefitId", claim.PatientBenefitID).Error("Patient benefit not found in DeleteClaim")
			return fiber.NewError(fiber.StatusNotFound, "Patient benefit not found")
		}
		uc.Log.WithError(err).Error("Failed to get patient benefit by ID in DeleteClaim")
		return err
	}
	if err := uc.PatientBenefitRepository.BalanceReduction(tx, patientBenefit, -(*claim.ApprovedAmount)); err != nil {
		uc.Log.WithError(err).Error("Failed to restore patient benefit balance in DeleteClaim")
		if err == gorm.ErrInvalidData {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid patient benefit data")
		}
		return err
	}

	if err := uc.Repository.Delete(tx, claim); err != nil {
		uc.Log.WithError(err).Error("Failed to delete claim")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction in DeleteClaim")
		return err
	}
	uc.Log.WithField("id", id).Info("Claim deleted successfully")
	return nil
}

func (uc *ClaimUseCase) GetAll(ctx context.Context, request *model.ClaimFilterQuery) ([]model.ClaimResponse, int64, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation error in GetAllClaims")
		return nil, 0, err
	}

	claims, total, err := uc.Repository.FindAllWithQuery(tx, request)
	if err != nil {
		uc.Log.WithError(err).Error("Error searching claims")
		return nil, 0, err
	}

	responses := make([]model.ClaimResponse, len(claims))
	for i, c := range claims {
		responses[i] = *converter.ClaimToResponse(&c)
	}
	return responses, total, nil
}