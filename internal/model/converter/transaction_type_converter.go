package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func TransactionTypeToResponse(transactionType *entity.TransactionType) *model.TransactionTypeResponse {
	return &model.TransactionTypeResponse{
		ID:   transactionType.ID,
		Name: transactionType.Name,
	}
}

