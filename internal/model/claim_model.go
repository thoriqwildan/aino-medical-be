package model

import (
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
)

type ClaimRequest struct {
	PatientID uint `json:"patient_id"`
	BenefitCode string `json:"benefit_code"`
	ClaimAmount float64 `json:"claim_amount"`
}

type PatientResponse struct {
	ID           uint      `json:"id"`
	Name        string    `json:"name"`
	BirthDate   helper.CustomDate    `json:"birth_date"`
	Gender      string 	`json:"gender"`
	PlanType    PlanTypeResponse `json:"plan_type"`
	Employee    *EmployeeResponse `json:"employee,omitempty"`
}

type ClaimResponse struct {
	ID 					uint    `json:"id"`
	ClaimAmount float64 `json:"claim_amount"`
	TransactionDate helper.CustomDate `json:"transaction_date"`
	SubmissionDate helper.CustomDate `json:"submission_date"`
	SLAStatus string `json:"sla_status"`
	ApprovedAmount float64 `json:"approved_amount"`
	ClaimStatus string `json:"claim_status"`
	MedicalFacility string `json:"medical_facility"`
	City string `json:"city"`
	Diagnosis string `json:"diagnosis"`
	DocLink string `json:"doc_link"`
	TransactionStatus string `json:"transaction_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	TransactionType TransactionTypeResponse `json:"transaction_type"`
	Patient PatientResponse `json:"patient"`
	Benefit BenefitResponse `json:"benefit"`
	Employee *EmployeeResponse `json:"employee,omitempty"`
}

type UpdateClaimRequest struct {
	ID                  uint      `json:"id" validate:"required"`
	ClaimAmount         float64   `json:"claim_amount" validate:"required"`
	TransactionTypeID   *uint     `json:"transaction_type_id"`
	TransactionDate     *helper.CustomDate `json:"transaction_date"`
	SubmissionDate      *helper.CustomDate `json:"submission_date"`
	SLA                 *string   `json:"sla" validate:"omitempty,oneof='meet' 'overdue'"`
	ClaimStatus         string    `json:"claim_status" validate:"required,oneof='On Plafond' 'Over Plafond'"`
	MedicalFacility     *string   `json:"medical_facility"`
	City                *string   `json:"city"`
	Diagnosis           *string   `json:"diagnosis"`
	DocLink             *string   `json:"doc_link"`
	TransactionStatus   string    `json:"transaction_status" validate:"required,oneof='Successful' 'Pending' 'Failed'"`
}

type ClaimFilterQuery struct {
  DateFrom          string                `form:"date_from"`
  DateTo            string                `form:"date_to"`
  Department        string                `form:"department"`
  TransactionType   string                `form:"transaction_type"`
  SLAStatus         entity.SLA            `form:"sla_status"`
  ClaimStatus       entity.ClaimStatus    `form:"claim_status"`
  TransactionStatus entity.TransactionStatus `form:"transaction_status"`
	Page int `json:"page,omitempty" validate:"omitempty,numeric"`
	Limit int `json:"limit,omitempty" validate:"omitempty,numeric"`
}