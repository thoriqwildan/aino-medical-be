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
	"github.com/thoriqwildan/aino-medical-be/pkg/constant"
	"gorm.io/gorm"
)

type ClaimUseCase struct {
	Repository                          *repository.ClaimRepository
	PatientBenefitRepository            *repository.PatientBenefitRepository
	TransactionTypeRepository           *repository.TransactionTypeRepository
	EmployeeRepository                  *repository.EmployeeRepository
	PatientYearlyBenefitClaimRepository *repository.PatientYearlyBenefitClaimRepository
	BenefitRepository                   *repository.BenefitRepository
	DepartmentRepository                *repository.DepartmentRepository
	PlanTypeRepository                  *repository.PlanTypeRepository
	YearlyBenefitClaimRepository        *repository.YearlyBenefitClaimRepository
	Log                                 *logrus.Logger
	DB                                  *gorm.DB
	Validate                            *validator.Validate
}

func NewClaimUseCase(
	repo *repository.ClaimRepository,
	db *gorm.DB,
	validate *validator.Validate,
	log *logrus.Logger,

	patientBenefitRepo *repository.PatientBenefitRepository,
	transactionTypeRepo *repository.TransactionTypeRepository,
	employeeRepo *repository.EmployeeRepository,
	patientYearlyBenefitClaimRepo *repository.PatientYearlyBenefitClaimRepository,
	benefitRepo *repository.BenefitRepository,
	departmentRepo *repository.DepartmentRepository,
	planTypeRepo *repository.PlanTypeRepository,
	yearlyBenefitClaimRepo *repository.YearlyBenefitClaimRepository,
) *ClaimUseCase {
	return &ClaimUseCase{
		Repository:                          repo,
		DB:                                  db,
		Validate:                            validate,
		Log:                                 log,
		PatientBenefitRepository:            patientBenefitRepo,
		TransactionTypeRepository:           transactionTypeRepo,
		EmployeeRepository:                  employeeRepo,
		PatientYearlyBenefitClaimRepository: patientYearlyBenefitClaimRepo,
		BenefitRepository:                   benefitRepo,
		DepartmentRepository:                departmentRepo,
		PlanTypeRepository:                  planTypeRepo,
		YearlyBenefitClaimRepository:        yearlyBenefitClaimRepo,
	}
}

func (uc ClaimUseCase) validateLimitationType(
	tx *gorm.DB,
	benefit *entity.Benefit,
	patientBenefit *entity.PatientBenefit,
	date time.Time,
) bool {
	if tx == nil || patientBenefit == nil {
		return false
	}

	q := tx.Model(&entity.Claim{}).
		Where("patient_benefit_id = ?", patientBenefit.ID).
		Where("transaction_status = ?", string(entity.TransactionStatusSuccessful))

	loc := date.Location()

	var start, end time.Time
	switch benefit.LimitationType {
	case entity.LimitationTypePerYear:
		start = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(1, 0, 0)
	case entity.LimitationTypePerMonth:
		start = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 1, 0)
	case entity.LimitationTypePerDay:
		start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		end = start.AddDate(0, 0, 1)
	case entity.LimitationTypePerPregnancy:
		start = date.AddDate(0, -9, 0)
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		end = dayStart.AddDate(0, 0, 1)
	default:
		return true
	}

	var cnt int64
	if err := q.Where("transaction_date >= ? AND transaction_date < ?", start, end).Count(&cnt).Error; err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(cnt)
	return cnt == 0
}

func (uc *ClaimUseCase) validatePatientRemainingAndYearly(
	tx *gorm.DB,
	amount float64,
	claimPatient entity.Patient,
	benefit *entity.Benefit,
	patientBenefit *entity.PatientBenefit,
) *fiber.Error {
	fmt.Println(patientBenefit)
	fmt.Println(amount)
	if patientBenefit != nil && patientBenefit.RemainingPlafond != nil {
		fmt.Println(*patientBenefit.RemainingPlafond)
		if benefit != nil &&
			benefit.YearlyBenefitClaimID != nil &&
			benefit.YearlyBenefitClaim != nil &&
			amount > *patientBenefit.RemainingPlafond {

			patientYBC, err := uc.PatientYearlyBenefitClaimRepository.FindOrCreate(tx, claimPatient, *benefit.YearlyBenefitClaim)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			if amount > patientYBC.YearlyClaimRemaining {
				return fiber.NewError(fiber.StatusBadRequest, "amount exceeds yearly remaining benefit claim limit")
			}
		}

		if amount > *patientBenefit.RemainingPlafond {
			return fiber.NewError(fiber.StatusBadRequest, "amount exceeds remaining plafond")
		}
		return nil
	}

	if benefit != nil && benefit.YearlyBenefitClaimID != nil && benefit.YearlyBenefitClaim != nil {
		patientYBC, err := uc.PatientYearlyBenefitClaimRepository.FindOrCreate(tx, claimPatient, *benefit.YearlyBenefitClaim)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if patientYBC != nil {
			if amount > patientYBC.YearlyClaimRemaining {
				return fiber.NewError(fiber.StatusBadRequest, "amount exceeds yearly remaining benefit claim limit")
			}
		}

	}
	return nil
}

func (uc *ClaimUseCase) validateBenefitPlafondAndYearly(
	amount float64,
	benefit *entity.Benefit,
) *fiber.Error {
	if benefit.Plafond != nil {
		if benefit.YearlyBenefitClaimID != nil && benefit.YearlyBenefitClaim != nil &&
			amount > *benefit.Plafond {

			if amount > benefit.YearlyBenefitClaim.YearlyClaim {
				return fiber.NewError(fiber.StatusBadRequest, "amount exceeds yearly benefit claim plafond limit")
			}
		}

		if amount > *benefit.Plafond {
			return fiber.NewError(fiber.StatusBadRequest, "amount exceeds benefit plafond")
		}
		return nil
	}

	if benefit.YearlyBenefitClaimID != nil && benefit.YearlyBenefitClaim != nil {
		fmt.Println(amount)
		fmt.Println(benefit.YearlyBenefitClaim.YearlyClaim)
		fmt.Println(amount > benefit.YearlyBenefitClaim.YearlyClaim)
		if amount > benefit.YearlyBenefitClaim.YearlyClaim {
			return fiber.NewError(fiber.StatusBadRequest, "amount exceeds yearly benefit claim plafond limit")
		}
	}
	return nil
}

func (uc *ClaimUseCase) Create(ctx context.Context, request *model.ClaimRequest) (*model.ClaimResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation error in ClaimRequest")
		return nil, err
	}

	now := time.Now()
	startDateOfCurrentYear := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	benefit := &entity.Benefit{}
	if err := uc.Repository.GetBenefitByCode(tx, benefit, request.BenefitCode); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.Log.Error("Benefit not found")
			return nil, fiber.NewError(fiber.StatusNotFound, "Benefit not found")
		}
		uc.Log.WithError(err).Error("Failed to get benefit by code")
		return nil, err
	}

	patient := &entity.Patient{}
	if err := uc.Repository.GetPatientByID(tx, patient, request.PatientID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	patientBenefit, err := uc.PatientBenefitRepository.FindOrCreate(tx, patient, benefit, benefit.Plafond, startDateOfCurrentYear, patient.Employee.ProRate)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find or create patient benefit")
		return nil, err
	}
	claim := &entity.Claim{
		PatientID:         request.PatientID,
		PatientBenefitID:  patientBenefit.ID,
		ClaimAmount:       request.ClaimAmount,
		TransactionStatus: entity.TransactionStatusPending,
	}

	errValidateRemaining := uc.validatePatientRemainingAndYearly(tx, request.ClaimAmount, *patient, benefit, patientBenefit)
	if errValidateRemaining != nil {
		if errValidateRemaining.Code == fiber.StatusBadRequest {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Claim "+errValidateRemaining.Error())
		}
		return nil, errValidateRemaining
	}

	// errValidatePlafond := uc.validateBenefitPlafondAndYearly(request.ClaimAmount, benefit)
	// if errValidatePlafond != nil {
	// 	if errValidatePlafond.Code == fiber.StatusBadRequest {
	// 		return nil, fiber.NewError(fiber.StatusBadRequest, "claim "+errValidatePlafond.Error())
	// 	}
	// 	return nil, errValidatePlafond
	// }

	isValidateLimitationType := uc.validateLimitationType(tx, benefit, patientBenefit, startDateOfCurrentYear)
	if !isValidateLimitationType {
		return nil, fiber.NewError(fiber.StatusBadRequest, "failed to create a claim due to type limitation validation")
	}

	if patientBenefit.InitialPlafond != nil {
		if *patientBenefit.InitialPlafond < request.ClaimAmount {
			claim.ClaimStatus = entity.ClaimStatusOverPlafond
		} else {
			claim.ClaimStatus = entity.ClaimStatusOnPlafond
		}
	} else {
		claim.ClaimStatus = entity.ClaimStatusOnPlafond
	}

	if patient.FamilyMemberID != nil {
		claim.EmployeeID = patient.FamilyMember.EmployeeID
	} else {
		claim.EmployeeID = *patient.EmployeeID
	}
	// DELETED: Cause claim amount not represent amount has been claim approve
	//if err := uc.PatientBenefitRepository.BalanceReduction(tx, patientBenefit, claim.ClaimAmount); err != nil {
	//	uc.Log.WithError(err).Error("Failed to reduce patient benefit balance")
	//	if errors.Is(err, gorm.ErrInvalidData) {
	//		return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient benefit balance")
	//	}
	//	return nil, err
	//}

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
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
			response.RemainingPlafond = rp
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
	startDateOfCurrentYear := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	claim := &entity.Claim{}
	if err := uc.Repository.GetByID(tx, claim, request.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.Log.WithField("id", request.ID).Error("Claim not found in UpdateClaim")
			return nil, fiber.NewError(fiber.StatusNotFound, "Claim not found")
		}
		uc.Log.WithError(err).Error("Failed to get claim by ID in UpdateClaim")
		return nil, err
	}

	benefit := &entity.Benefit{}
	if err := uc.BenefitRepository.GetById(tx, claim.PatientBenefit.BenefitID, benefit); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.Log.WithField("benefitId", claim.PatientBenefit.BenefitID).Error("Benefit not found in UpdateClaim")
			return nil, fiber.NewError(fiber.StatusNotFound, "Benefit not found")
		}
		uc.Log.WithError(err).Error("Failed to get benefit by ID in UpdateClaim")
		return nil, err
	}

	patientBenefit, err := uc.PatientBenefitRepository.FindOrCreate(tx, &claim.Patient, benefit, benefit.Plafond, startDateOfCurrentYear, claim.Employee.ProRate)
	if err != nil {
		uc.Log.WithError(err).Error("Failed to find or create patient benefit in UpdateClaim")
		return nil, err
	}

	errValidateRemaining := uc.validatePatientRemainingAndYearly(tx, request.ClaimAmount, claim.Patient, benefit, patientBenefit)
	if errValidateRemaining != nil {
		if errValidateRemaining.Code == fiber.StatusBadRequest {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Claim "+errValidateRemaining.Error())
		}
		return nil, errValidateRemaining
	}

	errValidatePlafond := uc.validateBenefitPlafondAndYearly(request.ClaimAmount, benefit)
	if errValidatePlafond != nil {
		if errValidatePlafond.Code == fiber.StatusBadRequest {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Claim "+errValidatePlafond.Error())
		}
		return nil, errValidatePlafond
	}

	if claim.ApprovedAmount != nil && patientBenefit.RemainingPlafond != nil {
		*patientBenefit.RemainingPlafond += *claim.ApprovedAmount
	}

	errValidateRemaining = uc.validatePatientRemainingAndYearly(tx, request.ApproveAmount, claim.Patient, benefit, patientBenefit)
	if errValidateRemaining != nil {
		if errValidateRemaining.Code == fiber.StatusBadRequest {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Approve "+errValidateRemaining.Error())
		}
		return nil, errValidateRemaining
	}

	errValidatePlafond = uc.validateBenefitPlafondAndYearly(request.ApproveAmount, benefit)
	if errValidatePlafond != nil {
		if errValidatePlafond.Code == fiber.StatusBadRequest {
			return nil, fiber.NewError(fiber.StatusBadRequest, "approve "+errValidatePlafond.Error())
		}
		return nil, errValidatePlafond
	}

	isValidateLimitationType := uc.validateLimitationType(tx, benefit, patientBenefit, startDateOfCurrentYear)
	if !isValidateLimitationType {
		return nil, fiber.NewError(fiber.StatusBadRequest, "failed to update a claim due to type limitation validation")
	}

	if (patientBenefit.RemainingPlafond != nil) && request.TransactionStatus == "Successful" {
		*patientBenefit.RemainingPlafond -= request.ApproveAmount
		tx.Save(&patientBenefit)
	}

	if (patientBenefit.RemainingPlafond == nil && benefit.YearlyBenefitClaim != nil) && request.TransactionStatus == "Successful" {
		patientYBC, err := uc.PatientYearlyBenefitClaimRepository.FindOrCreate(tx, claim.Patient, *benefit.YearlyBenefitClaim)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		patientYBC.YearlyClaimRemaining -= request.ApproveAmount
		tx.Save(&patientYBC)
	}

	claim.ClaimAmount = request.ClaimAmount
	claim.ApprovedAmount = &request.ApproveAmount
	claim.SLA = func(value *model.UpdateClaimRequest) *entity.SLA {
		if value.SLA != nil {
			slaStatus := helper.DetermineSLAStatus(*value.SLA)
			return &slaStatus
		} else {
			return claim.SLA
		}
	}(request)
	claim.TransactionTypeID = request.TransactionTypeID
	claim.TransactionStatus = entity.TransactionStatus(request.TransactionStatus)
	if patientBenefit.InitialPlafond != nil {
		if *patientBenefit.InitialPlafond < request.ClaimAmount {
			claim.ClaimStatus = entity.ClaimStatusOverPlafond
		} else {
			claim.ClaimStatus = entity.ClaimStatusOnPlafond
		}
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

func (uc ClaimUseCase) ImportClaims(ctx context.Context, request *model.ClaimsImportRequest) error {
	uc.Log.WithContext(ctx).WithField("claims_count", len(request.Claims)).
		Info("ImportClaims: start")

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		uc.Log.WithContext(ctx).Info("ImportClaims: defer rollback transaction")
		tx.Rollback()
	}()

	uc.Log.WithContext(ctx).Info("ImportClaims: validating request")
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithContext(ctx).WithError(err).
			Error("ImportClaims: validation failed")
		return err
	}
	uc.Log.WithContext(ctx).Info("ImportClaims: validation ok")

	now := time.Now()
	startDateOfCurrentYear := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	uc.Log.WithContext(ctx).WithFields(logrus.Fields{
		"now":                       now.Format(time.RFC3339),
		"startDateOfCurrentYearRaw": startDateOfCurrentYear.Format(time.RFC3339),
	}).Info("ImportClaims: computed dates")

	var claims []*entity.Claim

	for i, requestClaim := range request.Claims {
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"index":           i,
			"department":      requestClaim.Department,
			"employee_name":   requestClaim.EmployeeName,
			"patient_name":    requestClaim.PatientName,
			"plan_type_input": requestClaim.PlantType, // sesuai kode asli
			"benefit_code":    requestClaim.BenefitCode,
		}).Info("ImportClaims: processing claim item")

		uc.Log.WithContext(ctx).WithField("department", requestClaim.Department).
			Info("ImportClaims: FindOrCreate Department")
		deparment, errDepartment := uc.DepartmentRepository.FindOrCreate(tx, &entity.Department{Name: requestClaim.Department})
		if errDepartment != nil {
			uc.Log.WithContext(ctx).WithError(errDepartment).
				WithField("department", requestClaim.Department).
				Error("ImportClaims: Department FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errDepartment.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"department_id":   deparment.ID,
			"department_name": deparment.Name,
		}).Info("ImportClaims: Department ready")
		uc.Log.WithContext(ctx).WithField("plan_type_input", requestClaim.PlantType).
			Info("ImportClaims: FindOrCreate PlanType")
		planType, errPlanType := uc.PlanTypeRepository.FindOrCreate(tx, &entity.PlanType{Name: requestClaim.PlantType})
		if errPlanType != nil {
			uc.Log.WithContext(ctx).WithError(errPlanType).
				WithField("plan_type_input", requestClaim.PlantType).
				Error("ImportClaims: PlanType FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errPlanType.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"employee_name": requestClaim.EmployeeName,
			"bank_number":   requestClaim.BankNumber,
			"position":      deparment.Name, // sesuai kode asli
			"department_id": deparment.ID,
		}).Info("ImportClaims: FindOrCreate Employee")
		employee, errEmployee := uc.EmployeeRepository.FindOrCreate(tx, &entity.Employee{
			Name:         requestClaim.EmployeeName,
			Position:     deparment.Name, // sesuai kode asli
			Email:        helper.NameToGmail(requestClaim.EmployeeName),
			DepartmentID: deparment.ID,
			BankNumber:   requestClaim.BankNumber,
			Gender:       entity.GenderPreferNotSay,
			PlanTypeID:   planType.ID,
		})
		if errEmployee != nil {
			uc.Log.WithContext(ctx).WithError(errEmployee).
				WithField("employee_name", requestClaim.EmployeeName).
				Error("ImportClaims: Employee FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errEmployee.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"employee_id":      employee.ID,
			"employee_name":    employee.Name,
			"employee_prorate": employee.ProRate,
		}).Info("ImportClaims: Employee ready")

		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"patient_name": requestClaim.PatientName,
			"employee_id":  employee.ID,
		}).Info("ImportClaims: FindOrCreate Patient")
		patient, errPatient := uc.PatientBenefitRepository.FindOrCreatePatient(tx, &entity.Patient{
			Name:       requestClaim.PatientName,
			EmployeeID: &employee.ID,
			Gender:     entity.GenderPreferNotSay,
			PlanTypeID: employee.PlanTypeID,
		})
		if errPatient != nil {
			uc.Log.WithContext(ctx).WithError(errPatient).
				WithField("patient_name", requestClaim.PatientName).
				Error("ImportClaims: Patient FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errPatient.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"patient_id":   patient.ID,
			"patient_name": patient.Name,
		}).Info("ImportClaims: Patient ready")

		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"plan_type_id":   planType.ID,
			"plan_type_name": planType.Name,
		}).Info("ImportClaims: PlanType ready")

		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"benefit_code": requestClaim.BenefitCode,
			"plan_type":    requestClaim.PlantType,
		}).Info("ImportClaims: Build defaultBenefit via constant.DefaultBenefit")
		defaultBenefit := constant.DefaultBenefit(requestClaim.BenefitCode, requestClaim.PlantType)
		if defaultBenefit == nil {
			uc.Log.WithContext(ctx).WithFields(logrus.Fields{
				"benefit_code": requestClaim.BenefitCode,
				"plafond":      requestClaim.Plafond,
				"plan_type_id": planType.ID,
			}).Info("ImportClaims: defaultBenefit nil, using fallback struct")
			defaultBenefit = &entity.Benefit{
				Code:           requestClaim.BenefitCode,
				Plafond:        &requestClaim.Plafond,
				PlanTypeID:     planType.ID,
				LimitationType: entity.LimitationTypePerIncident,
			}
		}

		uc.Log.WithContext(ctx).WithField("benefit_code", defaultBenefit.Code).
			Info("ImportClaims: FindOrCreate Benefit")
		benefit, errBenefit := uc.BenefitRepository.FindOrCreate(tx, defaultBenefit)
		if errBenefit != nil {
			uc.Log.WithContext(ctx).WithError(errBenefit).
				WithField("benefit_code", defaultBenefit.Code).
				Error("ImportClaims: Benefit FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errBenefit.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"benefit_id":   benefit.ID,
			"benefit_code": benefit.Code,
		}).Info("ImportClaims: Benefit ready")

		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"patient_id":     patient.ID,
			"benefit_id":     benefit.ID,
			"plafond_input":  requestClaim.Plafond,
			"start_of_year?": startDateOfCurrentYear.Format(time.RFC3339),
			"employee_pr":    employee.ProRate,
		}).Info("ImportClaims: FindOrCreate PatientBenefit")
		patientBenefit, errPatientBenefit := uc.PatientBenefitRepository.FindOrCreate(
			tx, patient, benefit, &requestClaim.Plafond, startDateOfCurrentYear, employee.ProRate,
		)
		if errPatientBenefit != nil {
			uc.Log.WithContext(ctx).WithError(errPatientBenefit).
				WithFields(logrus.Fields{
					"patient_id": patient.ID,
					"benefit_id": benefit.ID,
				}).Error("ImportClaims: PatientBenefit FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errPatientBenefit.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"patient_benefit_id": patientBenefit.ID,
		}).Info("ImportClaims: PatientBenefit ready")

		uc.Log.WithContext(ctx).WithField("transaction_type_input", requestClaim.TransactionType).
			Info("ImportClaims: FindOrCreate TransactionType")
		transactionType, errTransactionType := uc.TransactionTypeRepository.FindOrCreate(
			tx, &entity.TransactionType{Name: *requestClaim.TransactionType},
		)
		if errTransactionType != nil {
			uc.Log.WithContext(ctx).WithError(errTransactionType).
				WithField("transaction_type_input", requestClaim.TransactionType).
				Error("ImportClaims: TransactionType FindOrCreate failed")
			return fiber.NewError(fiber.StatusInternalServerError, errTransactionType.Error())
		}
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"transaction_type_id":   transactionType.ID,
			"transaction_type_name": transactionType.Name,
		}).Info("ImportClaims: TransactionType ready")

		sla := helper.DetermineSLAStatus(requestClaim.SLA)
		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"sla_raw":     requestClaim.SLA,
			"sla_decoded": sla,
		}).Info("ImportClaims: SLA determined")

		uc.Log.WithContext(ctx).WithFields(logrus.Fields{
			"claim_amount":         requestClaim.ClaimAmount,
			"approved_amount":      requestClaim.ApproveAmount,
			"claim_status":         requestClaim.ClaimStatus,
			"transaction_status":   requestClaim.TransactionStatus,
			"transaction_date_raw": requestClaim.TransactionDate,
			"submission_date_raw":  requestClaim.SubmissionDate,
			"facility":             requestClaim.MedicalFacility,
			"city":                 requestClaim.City,
			"doc_link":             requestClaim.DocLink,
			"diagnosis":            requestClaim.Diagnosis,
		}).Info("ImportClaims: building Claim entity")

		claims = append(claims, &entity.Claim{
			PatientBenefitID:    patientBenefit.ID,
			PatientID:           patient.ID,
			EmployeeID:          employee.ID,
			ClaimAmount:         requestClaim.ClaimAmount,
			TransactionTypeID:   &transactionType.ID,
			TransactionDate:     (*time.Time)(requestClaim.TransactionDate),
			SubmissionDate:      (*time.Time)(requestClaim.SubmissionDate),
			SLA:                 &sla,
			ClaimStatus:         entity.ClaimStatus(requestClaim.ClaimStatus),
			ApprovedAmount:      requestClaim.ApproveAmount,
			MedicalFacilityName: &requestClaim.MedicalFacility,
			City:                &requestClaim.City,
			DocLink:             &requestClaim.DocLink,
			Diagnosis:           &requestClaim.Diagnosis,
			TransactionStatus:   entity.TransactionStatus(requestClaim.TransactionStatus),
		})
		uc.Log.WithContext(ctx).WithField("claims_len", len(claims)).
			Info("ImportClaims: claim appended")
	}

	uc.Log.WithContext(ctx).WithField("claims_len", len(claims)).
		Info("ImportClaims: UpsertBatches start")
	errUpsert := uc.Repository.UpsertBatches(tx, claims)
	if errUpsert != nil {
		uc.Log.WithContext(ctx).WithError(errUpsert).
			Error("ImportClaims: UpsertBatches failed")
		return fiber.NewError(fiber.StatusInternalServerError, errUpsert.Error())
	}
	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction in ImportClaims")
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	uc.Log.WithContext(ctx).Info("ImportClaims: UpsertBatches ok")

	uc.Log.WithContext(ctx).Info("ImportClaims: returning nil (no commit performed here)")
	return nil
}

func (uc *ClaimUseCase) GetClaim(ctx context.Context, id uint) (*model.ClaimResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claim := &entity.Claim{}
	if err := uc.Repository.GetByID(tx, claim, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.Log.WithField("id", id).Error("Claim not found in DeleteClaim")
			return fiber.NewError(fiber.StatusNotFound, "Claim not found")
		}
		uc.Log.WithError(err).Error("Failed to get claim by ID in DeleteClaim")
		return err
	}

	patientBenefit := &entity.PatientBenefit{}
	if err := uc.PatientBenefitRepository.FindById(tx, patientBenefit, claim.PatientBenefitID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			uc.Log.WithField("patientBenefitId", claim.PatientBenefitID).Error("Patient benefit not found in DeleteClaim")
			return fiber.NewError(fiber.StatusNotFound, "Patient benefit not found")
		}
		uc.Log.WithError(err).Error("Failed to get patient benefit by ID in DeleteClaim")
		return err
	}
	if claim.ApprovedAmount != nil {
		if err := uc.PatientBenefitRepository.BalanceReduction(tx, patientBenefit, -(*claim.ApprovedAmount)); err != nil {
			uc.Log.WithError(err).Error("Failed to restore patient benefit balance in DeleteClaim")
			if errors.Is(err, gorm.ErrInvalidData) {
				return fiber.NewError(fiber.StatusBadRequest, "Invalid patient benefit data")
			}
			return err
		}
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
