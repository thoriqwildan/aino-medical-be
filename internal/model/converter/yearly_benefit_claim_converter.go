package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func YearlyBenefitClaimToResponse(yearlyClaim *entity.YearlyBenefitClaim) *model.YearlyBenefitClaimResponse {
	yearlyBenefitClaim := &model.YearlyBenefitClaimResponse{
		ID:          yearlyClaim.ID,
		Code:        yearlyClaim.Code,
		YearlyClaim: yearlyClaim.YearlyClaim,
		CreatedAt:   yearlyClaim.CreatedAt,
		UpdatedAt:   yearlyClaim.UpdatedAt,
		Benefits:    []*model.BenefitResponse{},
	}
	if yearlyClaim.Benefits != nil || len(yearlyBenefitClaim.Benefits) > 0 {
		for _, benefit := range yearlyClaim.Benefits {
			yearlyBenefitClaim.Benefits = append(yearlyBenefitClaim.Benefits, BenefitToResponse(benefit))
		}
	}
	return yearlyBenefitClaim
}

func YearlyBenefitClaimToResponses(yearlyClaim []*entity.YearlyBenefitClaim) []*model.YearlyBenefitClaimResponse {
	yearlyClaims := make([]*model.YearlyBenefitClaimResponse, len(yearlyClaim))
	for i, yu := range yearlyClaim {
		yearlyClaims[i] = YearlyBenefitClaimToResponse(yu)
	}
	return yearlyClaims
}
