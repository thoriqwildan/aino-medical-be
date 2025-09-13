package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func BenefitToResponse(benefit *entity.Benefit) *model.BenefitResponse {
	response := &model.BenefitResponse{
		ID:             benefit.ID,
		Name:           benefit.Name,
		Detail:         benefit.Detail,
		Code:           benefit.Code,
		LimitationType: string(benefit.LimitationType),
	}
	if benefit.YearlyBenefitClaim != nil {
		response.YearlyBenefitClaim = YearlyBenefitClaimToResponse(benefit.YearlyBenefitClaim)
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

	return response
}
