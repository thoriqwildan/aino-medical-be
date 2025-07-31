package model

type ClaimRequest struct {
	PatientID string `json:"patient_id"`
	EmployeeID string `json:"employee_id"`
	BenefitCode string `json:"benefit_code"`
	ClaimAmount float64 `json:"claim_amount"`
}