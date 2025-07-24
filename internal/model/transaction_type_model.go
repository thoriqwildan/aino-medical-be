package model

type TransactionTypeRequest struct {
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type TransactionTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateTransactionTypeRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=3,max=255"`
}