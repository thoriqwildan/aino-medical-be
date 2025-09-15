package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func YearlyBenefitClaimToResponse(yearlyClaim *entity.YearlyBenefitClaim) *model.YearlyBenefitClaimResponse {
	return &model.YearlyBenefitClaimResponse{
		ID:          yearlyClaim.ID,
		Code:        yearlyClaim.Code,
		YearlyClaim: yearlyClaim.YearlyClaim,
		CreatedAt:   yearlyClaim.CreatedAt,
		UpdatedAt:   yearlyClaim.UpdatedAt,
	}
}

func YearlyBenefitClaimToResponses(yearlyClaim []*entity.YearlyBenefitClaim) []*model.YearlyBenefitClaimResponse {
	yearlyClaims := make([]*model.YearlyBenefitClaimResponse, len(yearlyClaim))
	for _, yu := range yearlyClaim {
		yearlyClaims = append(yearlyClaims, YearlyBenefitClaimToResponse(yu))
	}
	return yearlyClaims
}
