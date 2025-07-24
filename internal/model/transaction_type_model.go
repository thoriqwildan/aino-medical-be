package model

type TransactionTypeRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type TransactionTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}