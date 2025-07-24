package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func PlanTypeToResponse(planType *entity.PlanType) *model.PlanTypeResponse {
	return &model.PlanTypeResponse{
		ID:  planType.ID,
		Name: planType.Name,
		Description: planType.Description,
	}
}