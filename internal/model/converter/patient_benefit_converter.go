package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func PatientBenefitToResponse(patientBenefit *entity.PatientBenefit) *model.PatientBenefitResponse {
	response := &model.PatientBenefitResponse{
		ID:               patientBenefit.ID,
		PatientID:        patientBenefit.PatientID,
		BenefitID:        patientBenefit.BenefitID,
		RemainingPlafond: patientBenefit.RemainingPlafond,
		InitialPlafond:   patientBenefit.InitialPlafond,
		StartDate:        patientBenefit.StartDate,
		EndDate:          patientBenefit.EndDate,
		Status:           string(patientBenefit.Status),
		CreatedAt:        patientBenefit.CreatedAt,
		UpdatedAt:        patientBenefit.UpdatedAt,
	}
	if patientBenefit.Patient.ID != 0 {
		response.Patient = PatientToResponse(&patientBenefit.Patient)
	}
	if patientBenefit.Benefit.ID != 0 {
		response.Benefit = BenefitToResponse(&patientBenefit.Benefit)
	}
	if len(patientBenefit.Claims) > 0 || patientBenefit.Claims != nil {
		for _, claim := range patientBenefit.Claims {
			response.Claims = append(response.Claims, ClaimToResponse(&claim))
		}
	}
	return response
}
