package converter

import (
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func EmployeeToResponse(employee *entity.Employee) *model.EmployeeResponse {
	response := &model.EmployeeResponse{
		ID:          employee.ID,
		Name:        employee.Name,
		Email:       employee.Email,
		Phone:       employee.Phone,
		Position:    employee.Position,
		BirthDate:   helper.CustomDate(employee.BirthDate),
		Gender:      string(employee.Gender),
		Dependences: *employee.Dependence,
		BankNumber:  employee.BankNumber,
		JoinDate:    helper.CustomDate(employee.JoinDate),
	}

	if len(employee.FamilyMembers) > 0 {
		for _, familyMember := range employee.FamilyMembers {
			response.FamilyMembers = append(response.FamilyMembers, *FamilyMemberToResponse(&familyMember))
		}
	}

	if employee.PlanType.ID != 0 || employee.PlanType.Name != "" {
		response.PlanType = *PlanTypeToResponse(&employee.PlanType)
	}

	if employee.Department.ID != 0 || employee.Department.Name != "" {
		response.Department = *DepartmentToResponse(&employee.Department)
	}

	if len(employee.FamilyMembers) > 0 {
		response.FamilyMembers = make([]model.FamilyMemberResponse, len(employee.FamilyMembers))
		for i, familyMember := range employee.FamilyMembers {
			response.FamilyMembers[i] = *FamilyMemberToResponse(&familyMember)
		}
	}

	return response
}
