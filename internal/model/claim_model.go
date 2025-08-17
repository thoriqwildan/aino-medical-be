package model

import (
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
)

type ClaimRequest struct {
	PatientID   uint    `json:"patient_id" validate:"required"`
	BenefitCode string  `json:"benefit_code" validate:"required"`
	ClaimAmount float64 `json:"claim_amount" validate:"required"`
}

type PatientResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	BirthDate helper.CustomDate `json:"birth_date"`
	Gender    string            `json:"gender"`
	PlanType  PlanTypeResponse  `json:"plan_type"`
	Employee  *EmployeeResponse `json:"employee,omitempty"`
}

type ClaimResponse struct {
	ID                uint                    `json:"id"`
	ClaimAmount       float64                 `json:"claim_amount"`
	TransactionDate   helper.CustomDate       `json:"transaction_date"`
	SubmissionDate    helper.CustomDate       `json:"submission_date"`
	SLAStatus         string                  `json:"sla_status"`
	ApprovedAmount    float64                 `json:"approved_amount"`
	ClaimStatus       string                  `json:"claim_status"`
	MedicalFacility   string                  `json:"medical_facility"`
	City              string                  `json:"city"`
	Diagnosis         string                  `json:"diagnosis"`
	DocLink           string                  `json:"doc_link"`
	TransactionStatus string                  `json:"transaction_status"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         *time.Time              `json:"updated_at"`
	TransactionType   TransactionTypeResponse `json:"transaction_type"`
	Patient           PatientResponse         `json:"patient"`
	Benefit           BenefitResponse         `json:"benefit"`
	Employee          *EmployeeResponse       `json:"employee,omitempty"`
}

type UpdateClaimRequest struct {
	ID                uint               `json:"id" validate:"required"`
	ClaimAmount       float64            `json:"claim_amount" validate:"required"`
	TransactionTypeID *uint              `json:"transaction_type_id" validate:"required"`
	TransactionDate   *helper.CustomDate `json:"transaction_date" validate:"required"`
	SubmissionDate    *helper.CustomDate `json:"submission_date" validate:"required"`
	SLA               *string            `json:"sla" validate:"required,oneof='meet' 'overdue'"`
	ClaimStatus       string             `json:"claim_status" validate:"required,oneof='On Plafond' 'Over Plafond'"`
	MedicalFacility   *string            `json:"medical_facility" validate:"required"`
	City              *string            `json:"city" validate:"required"`
	Diagnosis         *string            `json:"diagnosis" validate:"required"`
	DocLink           *string            `json:"doc_link" validate:"required"`
	TransactionStatus string             `json:"transaction_status" validate:"required,oneof='Successful' 'Pending' 'Failed'"`
}

type ClaimFilterQuery struct {
	SearchValue       string                   `query:"search_value" form:"search_value"`
	BenefitID         string                   `query:"benefit_id" form:"benefit_id"`
	RelationshipType  string                   `query:"relationship_type" form:"relationship_type" validate:"omitempty,oneof=wife husband father mother child"`
	DateFrom          string                   `form:"date_from" query:"date_from"`
	DateTo            string                   `form:"date_to" query:"date_to"`
	Department        string                   `form:"department" query:"department"`
	TransactionType   string                   `form:"transaction_type" query:"transaction_type"`
	SLAStatus         entity.SLA               `form:"sla_status" query:"sla_status"`
	ClaimStatus       entity.ClaimStatus       `form:"claim_status" query:"claim_status"`
	TransactionStatus entity.TransactionStatus `form:"transaction_status" query:"transaction_status"`
	Page              int                      `json:"page,omitempty" query:"page,omitempty"  validate:"omitempty,numeric"`
	Limit             int                      `json:"limit,omitempty" query:"limit,omitempty" validate:"omitempty,numeric"`
}
