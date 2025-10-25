package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func PatientYearlyBenefitClaimToResponse(claim *entity.PatientYearlyBenefitClaim) *model.PatientYearlyBenefitClaimResponse {
	response := &model.PatientYearlyBenefitClaimResponse{
		ID:                   claim.ID,
		PatientID:            claim.PatientID,
		YearlyBenefitClaimID: claim.YearlyBenefitClaimID,
		YearlyRemaining:      claim.YearlyClaimRemaining,
		CreatedAt:            claim.CreatedAt,
		UpdatedAt:            claim.UpdatedAt,
	}
	if claim.YearlyBenefitClaim.ID != 0 {
		response.YearlyClaimBenefit = YearlyBenefitClaimToResponse(&claim.YearlyBenefitClaim)
	}
	if claim.Patient.ID != 0 {
		response.Patient = PatientToResponse(&claim.Patient)
	}
	return response
}

func PatientYearlyBenefitClaimToResponses(claims []entity.PatientYearlyBenefitClaim) []*model.PatientYearlyBenefitClaimResponse {
	responses := make([]*model.PatientYearlyBenefitClaimResponse, len(claims))
	for i, claim := range claims {
		responses[i] = PatientYearlyBenefitClaimToResponse(&claim)
	}
	return responses
}
