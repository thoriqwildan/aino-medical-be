package model

import (
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
)

type EmployeeRequest struct {
	Name		 string `json:"name" validate:"required"`
	DepartmentID uint   `json:"department_id" validate:"required"`
	Position		 string `json:"position" validate:"required"`
	Email		 string `json:"email" validate:"required,email"`
	Phone		 string `json:"phone" validate:"required"`
	BirthDate	 helper.CustomDate `json:"birth_date" validate:"required"`
	Gender 	 string `json:"gender" validate:"required,oneof=male female"`
	PlanTypeID	 uint   `json:"plan_type_id" validate:"required"`
	Dependences string `json:"dependence,omitempty" validate:"omitempty"`
	BankNumber	 string `json:"bank_number" validate:"required"`
	JoinDate	 helper.CustomDate `json:"join_date" validate:"required"`
}

type EmployeeResponse struct {
	ID            uint   `json:"id"`
	Name		 			string `json:"name"`
	Position		 	string `json:"position"`
	Email		 			string `json:"email"`
	Phone					string `json:"phone"`
	BirthDate	 		helper.CustomDate `json:"birth_date"`
	Gender 	 			string `json:"gender"`
	Dependences 	string `json:"dependence,omitempty"`
	BankNumber	 	string `json:"bank_number"`
	JoinDate	 		helper.CustomDate `json:"join_date"`
	PlanType	 		PlanTypeResponse   `json:"plan_type"`
	Department   	DepartmentResponse   `json:"department"`
}

type UpdateEmployeeRequest struct {
	ID            uint   `json:"id" validate:"required"`
	Name		 string `json:"name" validate:"required"`
	DepartmentID uint   `json:"department_id" validate:"required"`
	Position		 string `json:"position" validate:"required"`
	Email		 string `json:"email" validate:"required,email"`
	Phone		 string `json:"phone" validate:"required"`
	BirthDate	 helper.CustomDate `json:"birth_date" validate:"required"`
	Gender 	 string `json:"gender" validate:"required,oneof=male female"`
	PlanTypeID	 uint   `json:"plan_type_id" validate:"required"`
	Dependences string `json:"dependence,omitempty" validate:"omitempty"`
	BankNumber	 string `json:"bank_number" validate:"required"`
	JoinDate	 helper.CustomDate `json:"join_date" validate:"required"`
}