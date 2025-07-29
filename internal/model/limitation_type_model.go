package model

type LimitationTypeRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
}

type LimitationTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateLimitationTypeRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1,max=100"`
}