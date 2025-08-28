package converter

import (
	"fmt"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func BenefitToResponse(benefit *entity.Benefit) *model.BenefitResponse {
	response := &model.BenefitResponse{
		ID:     benefit.ID,
		Name:   benefit.Name,
		Detail: benefit.Detail,
		Code:   benefit.Code,
	}

	if benefit.YearlyMax != nil {
		response.YearlyMax = benefit.YearlyMax
	} else {
		response.YearlyMax = nil
	}
	if benefit.Plafond != nil {
		response.Plafond = benefit.Plafond
	} else {
		benefit.Plafond = nil
	}

	if benefit.PlanType.ID != 0 {
		response.PlanType = *PlanTypeToResponse(&benefit.PlanType)
	}
	if benefit.LimitationType.ID != 0 {
		response.LimitationType = *LimitationTypeToResponse(&benefit.LimitationType)
	}
	fmt.Println(response)
	return response
}
