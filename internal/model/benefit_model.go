package model

type CreateBenefitRequest struct {
	Name           string   `json:"name" validate:"required,min=3,max=255"`
	PlanTypeID     uint     `json:"plan_type_id" validate:"required"`
	YearlyClaimID  *uint    `json:"yearly_claim_id" validate:"omitempty"`
	Detail         *string  `json:"detail,omitempty" validate:"omitempty,max=500"`
	Code           string   `json:"code" validate:"required,min=3,max=50"`
	LimitationType string   `json:"limitation_type" validate:"required,oneof='Per Year' 'Per Month' 'Per Day' 'Per Incident' 'Per Pregnancy'"`
	Plafond        float64  `json:"plafond,omitempty" validate:"omitempty,numeric"`
	YearlyMax      *float64 `json:"yearly_max" validate:"omitempty,numeric"`
}

type BenefitResponse struct {
	ID                 uint                        `json:"id"`
	Name               string                      `json:"name"`
	Detail             *string                     `json:"detail,omitempty"`
	Code               string                      `json:"code"`
	Plafond            *float64                    `json:"plafond"`
	YearlyMax          *float64                    `json:"yearly_max"`
	RemainingPlafond   *float64                    `json:"remaining_plafond,omitempty"`
	PlanType           PlanTypeResponse            `json:"plan_type"`
	YearlyBenefitClaim *YearlyBenefitClaimResponse `json:"yearly_benefit_claim,omitempty"`
	LimitationType     string                      `json:"limitation_type"`
}

type UpdateBenefitRequest struct {
	ID             uint    `json:"id" validate:"required"`
	Name           string  `json:"name" validate:"required,min=3,max=255"`
	PlanTypeID     uint    `json:"plan_type_id" validate:"required"`
	YearlyClaimID  *uint   `json:"yearly_claim_id" validate:"omitempty"`
	Detail         *string `json:"detail,omitempty" validate:"omitempty,max=500"`
	Code           string  `json:"code" validate:"required,min=3,max=50"`
	LimitationType string  `json:"limitation_type" validate:"required,oneof='Per Year' 'Per Month' 'Per Day' 'Per Incident' 'Per Pregnancy'"`
	Plafond        float64 `json:"plafond,omitempty" validate:"omitempty,numeric"`
	YearlyMax      float64 `json:"yearly_max,omitempty" validate:"omitempty,numeric"`
}
