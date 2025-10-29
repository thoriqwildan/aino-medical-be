package model

import "time"

type PatientYearlyBenefitClaimParams struct {
	PatientID            uint `validate:"required,gt=0"`
	YearlyBenefitClaimID uint `validate:"required,gt=0"`
}

type PatientYearlyBenefitClaimResponse struct {
	ID                   uint                        `json:"id"`
	PatientID            uint                        `json:"patient_id"`
	YearlyBenefitClaimID uint                        `json:"yearly_benefit_claim_id"`
	YearlyRemaining      float64                     `json:"yearly_remaining"`
	Patient              *PatientResponse            `json:"patient,omitempty"`
	YearlyClaimBenefit   *YearlyBenefitClaimResponse `json:"benefit,omitempty"`
	CreatedAt            time.Time                   `json:"created_at"`
	UpdatedAt            time.Time                   `json:"updated_at"`
}
