package model

type CreateBenefitRequest struct {
	Name             string  `json:"name" validate:"required,min=3,max=255"`
	PlanTypeID       uint    `json:"plan_type_id" validate:"required"`
	Detail           *string `json:"detail,omitempty" validate:"omitempty,max=500"`
	Code             string  `json:"code" validate:"required,min=3,max=50"`
	LimitationTypeID uint    `json:"limitation_type_id" validate:"required"`
	Plafond          float64 `json:"plafond,omitempty" validate:"omitempty,numeric"`
	YearlyMax        float64 `json:"yearly_max,omitempty" validate:"omitempty,numeric"`
}

type BenefitResponse struct {
	ID               uint                   `json:"id"`
	Name             string                 `json:"name"`
	Detail           *string                `json:"detail,omitempty"`
	Code             string                 `json:"code"`
	Plafond          *float64               `json:"plafond"`
	YearlyMax        *float64               `json:"yearly_max"`
	RemainingPlafond *float64               `json:"remaining_plafond,omitempty"`
	PlanType         PlanTypeResponse       `json:"plan_type"`
	LimitationType   LimitationTypeResponse `json:"limitation_type"`
}

type UpdateBenefitRequest struct {
	ID               uint    `json:"id" validate:"required"`
	Name             string  `json:"name" validate:"required,min=3,max=255"`
	PlanTypeID       uint    `json:"plan_type_id" validate:"required"`
	Detail           *string `json:"detail,omitempty" validate:"omitempty,max=500"`
	Code             string  `json:"code" validate:"required,min=3,max=50"`
	LimitationTypeID uint    `json:"limitation_type_id" validate:"required"`
	Plafond          float64 `json:"plafond,omitempty" validate:"omitempty,numeric"`
	YearlyMax        float64 `json:"yearly_max,omitempty" validate:"omitempty,numeric"`
}
