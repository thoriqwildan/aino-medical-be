package model

type PlanTypeRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=1"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type PlanTypeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description *string `json:"description,omitempty"`
}