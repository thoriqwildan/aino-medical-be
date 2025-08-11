package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func BenefitToResponse(benefit *entity.Benefit) *model.BenefitResponse {
  response := &model.BenefitResponse{
    ID:          benefit.ID,
    Name:        benefit.Name,
    Detail:      benefit.Detail,
    Code:        benefit.Code,
    Plafond:     &benefit.Plafond,
    YearlyMax:   &benefit.YearlyMax,
  }

  if benefit.PlanType.ID != 0 {
    response.PlanType = *PlanTypeToResponse(&benefit.PlanType)
  }
  if benefit.LimitationType.ID != 0 {
    response.LimitationType = *LimitationTypeToResponse(&benefit.LimitationType)
  }
  return response
}