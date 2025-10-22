package converter

import (
	"time" // Pastikan time diimport untuk time.Time{}
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func ClaimToResponse(claim *entity.Claim) *model.ClaimResponse {
	var transactionDate helper.CustomDate
	var submissionDate helper.CustomDate
	var slaStatus string
	var medicalFacility string
	var city string
	var diagnosis string
	var docLink string
	var updatedAt *time.Time

	if claim.TransactionDate != nil {
		transactionDate = helper.CustomDate(*claim.TransactionDate)
	} else {
		transactionDate = helper.CustomDate(time.Time{})
	}

	if claim.SubmissionDate != nil {
		submissionDate = helper.CustomDate(*claim.SubmissionDate)
	} else {
		submissionDate = helper.CustomDate(time.Time{})
	}

	if claim.SLA != nil {
		slaStatus = string(*claim.SLA)
	} else {
		slaStatus = ""
	}

	if claim.MedicalFacilityName != nil {
		medicalFacility = *claim.MedicalFacilityName
	} else {
		medicalFacility = ""
	}

	if claim.City != nil {
		city = *claim.City
	} else {
		city = ""
	}

	if claim.Diagnosis != nil {
		diagnosis = *claim.Diagnosis
	} else {
		diagnosis = ""
	}

	if claim.DocLink != nil {
		docLink = *claim.DocLink
	} else {
		docLink = ""
	}

	if claim.UpdatedAt != nil {
		updatedAt = claim.UpdatedAt
	} else {
		updatedAt = nil
	}

	result := &model.ClaimResponse{
		ID:                claim.ID,
		ClaimAmount:       claim.ClaimAmount,
		TransactionDate:   transactionDate,
		SubmissionDate:    submissionDate,
		SLAStatus:         slaStatus,
		ApprovedAmount:    claim.ApprovedAmount,
		ClaimStatus:       string(claim.ClaimStatus),
		MedicalFacility:   medicalFacility,
		City:              city,
		Diagnosis:         diagnosis,
		DocLink:           docLink,
		TransactionStatus: string(claim.TransactionStatus),
		CreatedAt:         claim.CreatedAt,
		UpdatedAt:         updatedAt,
	}

	if claim.TransactionType != nil {
		result.TransactionType = TransactionTypeToResponse(claim.TransactionType)
	}

	if claim.Patient.ID != 0 {
		result.Patient = &model.PatientResponse{
			ID:                                claim.Patient.ID,
			Name:                              claim.Patient.Name,
			BirthDate:                         helper.CustomDate(claim.Patient.BirthDate),
			Gender:                            string(claim.Patient.Gender),
			PatientYearlyBenefitClaimResponse: []*model.PatientYearlyBenefitClaimResponse{},
		}
	}
	if claim.Patient.PatientYearlyBenefitClaims != nil || len(claim.Patient.PatientYearlyBenefitClaims) > 0 {
		result.Patient.PatientYearlyBenefitClaimResponse = PatientYearlyBenefitClaimToResponses(claim.Patient.PatientYearlyBenefitClaims)
	}

	if claim.Employee.ID != 0 {
		result.Employee = EmployeeToResponse(&claim.Employee)
	}

	if claim.PatientBenefit.BenefitID != 0 {
		result.Benefit = BenefitToResponse(&claim.PatientBenefit.Benefit)
	}

	return result
}

func PatientToResponse(patient *entity.Patient) *model.PatientResponse {
	result := &model.PatientResponse{
		ID:        patient.ID,
		Name:      patient.Name,
		BirthDate: helper.CustomDate(patient.BirthDate),
		Gender:    string(patient.Gender),
		PlanType:  *PlanTypeToResponse(&patient.PlanType),
	}

	if patient.Employee != nil {
		result.Employee = EmployeeToResponse(patient.Employee)
	} else if patient.FamilyMember != nil && patient.FamilyMember.Employee != nil {
		result.Employee = EmployeeToResponse(patient.FamilyMember.Employee)
	} else {
		result.Employee = nil
	}
	return result
}

// func ClaimImportRequestToEntity(claim *model.ClaimImportRequest) *entity.Claim {
// 	claimEntity := &entity.Claim{
// 		ID:                claim.ID,
// 		PatientBenefitID:  claim.PatientID,
// 		PatientID:         claim.PatientBenefitID,
// 		EmployeeID:        claim.EmployeeID,
// 		ClaimAmount:       claim.ClaimAmount,
// 		TransactionTypeID: claim.TransactionTypeID,

// 		ApprovedAmount:      claim.ApproveAmount,
// 		ClaimStatus:         entity.ClaimStatus(claim.ClaimStatus),
// 		MedicalFacilityName: claim.MedicalFacility,
// 		City:                claim.City,
// 		Diagnosis:           claim.Diagnosis,
// 		DocLink:             claim.DocLink,
// 		TransactionStatus:   entity.TransactionStatus(claim.TransactionStatus),
// 	}
// 	if claim.TransactionDate != nil {
// 		claimEntity.TransactionDate = (*time.Time)(claim.TransactionDate)
// 	}
// 	if claim.SubmissionDate != nil {
// 		claimEntity.SubmissionDate = (*time.Time)(claim.SubmissionDate)
// 	}
// 	if claim.SLA != nil {
// 		claimEntity.SLA = func(value *model.ClaimImportRequest) *entity.SLA {
// 			if value.SLA != nil {
// 				slaStatus := helper.DetermineSLAStatus(*value.SLA)
// 				return &slaStatus
// 			} else {
// 				return nil
// 			}
// 		}(claim)
// 	}

// 	return claimEntity
// }

// func ClaimsImportRequestToEntity(claims *model.ClaimsImportRequest) []*entity.Claim {
// 	var claimsEntity []*entity.Claim
// 	for _, claim := range claims.Claims {
// 		claimsEntity = append(claimsEntity, ClaimImportRequestToEntity(claim))
// 	}
// 	return claimsEntity
// }
