package entity

type Genders string

const (
	GenderMale   Genders = "male"
	GenderFemale Genders = "female"
)

type SLA string

const (
	SLAMeet    SLA = "meet"
	SLAOverdue SLA = "overdue"
)

type ClaimStatus string

const (
	ClaimStatusOnPlafond   ClaimStatus = "On Plafond"
	ClaimStatusOverPlafond ClaimStatus = "Over Plafond"
)

type TransactionStatus string

const (
	TransactionStatusSuccessful TransactionStatus = "Successful"
	TransactionStatusPending    TransactionStatus = "Pending"
	TransactionStatusFailed     TransactionStatus = "Failed"
)

type PatientBenefitStatus string

const (
	PatientBenefitStatusActive    PatientBenefitStatus = "active"
	PatientBenefitStatusExhausted PatientBenefitStatus = "exhausted"
	PatientBenefitStatusExpired   PatientBenefitStatus = "expired"
)