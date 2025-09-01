package model

import "github.com/thoriqwildan/aino-medical-be/internal/helper"

type FamilyMemberRequest struct {
	Name             string            `json:"name" validate:"required"`
	EmployeeID       uint              `json:"employee_id" validate:"required"`
	BirthDate        helper.CustomDate `json:"birth_date" validate:"required"`
	RelationshipType string            `json:"relationship_type" validate:"required"`
	Gender           string            `json:"gender" validate:"required,oneof=male female"`
}

type FamilyMemberResponse struct {
	ID               uint              `json:"id"`
	Name             string            `json:"name" validate:"required"`
	RelationshipType string            `json:"relationship_type" validate:"required"`
	BirthDate        helper.CustomDate `json:"birth_date" validate:"required"`
	Gender           string            `json:"gender" validate:"required,oneof=male female"`
	PlanType         PlanTypeResponse  `json:"plan_type,omitempty" validate:"required"`
	Employee         EmployeeResponse  `json:"employee,omitempty" validate:"required"`
}

type UpdateFamilyMemberRequest struct {
	ID        *uint             `json:"id,omitempty" validate:"required,omitempty"`
	Name      string            `json:"name" validate:"required"`
	BirthDate helper.CustomDate `json:"birth_date" validate:"required"`
	Gender    string            `json:"gender" validate:"required,oneof=male female"`
}
