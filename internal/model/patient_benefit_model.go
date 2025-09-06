package model

import (
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/helper"
)

type UpdatePatientBenefitRequest struct {
	StartDate        *helper.CustomDate `json:"start_date" validate:"required"`
	EndDate          *helper.CustomDate `json:"end_date" validate:"required"`
	RemainingPlafond *float64           `json:"remaining_plafond" validate:"omitempty,gte=0"`
	InitialPlafond   *float64           `json:"initial_plafond" validate:"omitempty,gte=0"`
	YearlyMax        *float64           `json:"yearly_max" validate:"omitempty,gte=0"`
	Status           string             `json:"status" validate:"required,oneof='active' 'exhausted' 'expired'"`
}

type PatientBenefitParams struct {
	PatientID uint `validate:"required,gt=0"`
	BenefitID uint `validate:"required,gt=0"`
}

type PatientBenefitResponse struct {
	ID               uint             `json:"id"`
	PatientID        uint             `json:"patient_id"`
	BenefitID        uint             `json:"benefit_id"`
	RemainingPlafond *float64         `json:"remaining_plafond"`
	InitialPlafond   *float64         `json:"initial_plafond"`
	YearlyMax        *float64         `json:"yearly_max"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          *time.Time       `json:"end_date"`
	Status           string           `json:"status"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        *time.Time       `json:"updated_at"`
	Patient          *PatientResponse `json:"patient,omitempty"`
	Benefit          *BenefitResponse `json:"benefit,omitempty"`
	Claims           []*ClaimResponse `json:"claims,omitempty"`
}
