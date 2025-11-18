package constant

import "github.com/thoriqwildan/aino-medical-be/internal/entity"

var DEFAULT_YEARLY_CLAIM = []*entity.YearlyBenefitClaim{
	{
		Code:        "RI-PLANB",
		YearlyClaim: 75000000,
	},
	{
		Code:        "RI-PLANC",
		YearlyClaim: 60000000,
	},
	{
		Code:        "RI-PLAND",
		YearlyClaim: 50000000,
	},
	{
		Code:        "RJ-PLANA",
		YearlyClaim: 12000000,
	},
	{
		Code:        "RJ-PLANB",
		YearlyClaim: 9000000,
	},
	{
		Code:        "RJ-PLANC",
		YearlyClaim: 8000000,
	},
	{
		Code:        "RJ-PLAND",
		YearlyClaim: 6000000,
	},
	{
		Code:        "RG-PLANA",
		YearlyClaim: 5000000,
	},
	{
		Code:        "RG-PLANB",
		YearlyClaim: 4000000,
	},
	{
		Code:        "RG-PLANC",
		YearlyClaim: 3000000,
	},
	{
		Code:        "RG-PLAND",
		YearlyClaim: 2500000,
	},
}

func DefaultYearlyBenefitClaim(code string) *entity.YearlyBenefitClaim {
	for _, rw := range DEFAULT_YEARLY_CLAIM {
		if rw.Code == code {
			return rw
		}
	}
	return nil
}