package model

import "time"

type YearlyBenefitClaimRequest struct {
	Code        string  `validate:"required,max=255" json:"code"`
	YearlyClaim float64 `validate:"required,numeric" json:"yearly_claim"`
}

type UpdateYearlyBenefitClaimRequest struct {
	Code        string  `validate:"required,max=255" json:"code"`
	YearlyClaim float64 `validate:"required,numeric" json:"yearly_claim"`
}

type YearlyBenefitClaimResponse struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	YearlyClaim float64   `json:"yearly_claim"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type YearlyBenefitClaimFilter struct {
	Code  string `json:"code" query:"code" validate:"omitempty,max=255"`
	Page  int    `json:"page,omitempty" query:"page,omitempty"  validate:"omitempty,numeric"`
	Limit int    `json:"limit,omitempty" query:"limit,omitempty" validate:"omitempty,numeric"`
}
