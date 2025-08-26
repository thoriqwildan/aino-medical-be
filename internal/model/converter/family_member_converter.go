package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func FamilyMemberToResponse(familyMember *entity.FamilyMember) *model.FamilyMemberResponse {
	result := &model.FamilyMemberResponse{
		ID:               familyMember.ID,
		Name:             familyMember.Name,
		BirthDate:        helper.CustomDate(familyMember.BirthDate),
		Gender:           string(familyMember.Gender),
		RelationshipType: string(familyMember.RelationshipType),
	}

	if familyMember.PlanType.ID != 0 || familyMember.PlanType.Name != "" {
		result.PlanType = *PlanTypeToResponse(&familyMember.PlanType)
	}

	if familyMember.Employee != nil {
		result.Employee = *EmployeeToResponse(familyMember.Employee)
	}

	return result
}
