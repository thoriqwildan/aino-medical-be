package entity

type Genders string

const (
	GenderMale   Genders = "male"
	GenderFemale Genders = "female"
)

// SLA Enum
type SLA string

const (
	SLAMeet    SLA = "meet"
	SLAOverdue SLA = "overdue"
)

// ClaimStatus Enum
type ClaimStatus string

const (
	ClaimStatusOnPlafond   ClaimStatus = "On Plafond"
	ClaimStatusOverPlafond ClaimStatus = "Over Plafond"
)

// TransactionStatus Enum
type TransactionStatus string

const (
	TransactionStatusSuccessful TransactionStatus = "Successful"
	TransactionStatusPending    TransactionStatus = "Pending"
	TransactionStatusFailed     TransactionStatus = "Failed"
)