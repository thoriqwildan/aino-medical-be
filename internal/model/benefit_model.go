package model

type CreateBenefitRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
	PlanTypeID uint `json:"plan_type_id" validate:"required"`
	Detail *string `json:"detail,omitempty" validate:"omitempty,max=500"`
	Code string `json:"code" validate:"required,min=3,max=50"`
	LimitationTypeID uint `json:"limitation_type_id" validate:"required"`
	Plafond float64 `json:"plafond,omitempty" validate:"omitempty,numeric"`
	YearlyMax float64 `json:"yearly_max,omitempty" validate:"omitempty,numeric"`
}