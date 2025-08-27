package entity

import (
	"time"

	"gorm.io/gorm"
)

type PlanType struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"unique;not null"`
	Description   *string        // TEXT bisa diwakili *string
	Benefits      []Benefit      `gorm:"foreignKey:PlanTypeID"`
	Employees     []Employee     `gorm:"foreignKey:PlanTypeID"`
	FamilyMembers []FamilyMember `gorm:"foreignKey:PlanTypeID"`
}

type TransactionType struct {
	ID     uint    `gorm:"primaryKey;autoIncrement"`
	Name   string  `gorm:"unique;not null"`
	Claims []Claim `gorm:"foreignKey:TransactionTypeID"`
}

type Department struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	Name      string          `gorm:"unique;not null"`
	CreatedAt time.Time       `gorm:"not null;autoCreateTime"`
	UpdatedAt *time.Time      `gorm:"autoUpdateTime"` // NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	Employees []Employee      `gorm:"foreignKey:DepartmentID"`
}

type Employee struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"not null"`
	DepartmentID  uint           `gorm:"not null"`
	Position      string         `gorm:"not null"`
	Email         string         `gorm:"unique;not null"`
	Phone         string         `gorm:"not null"`
	BirthDate     time.Time      `gorm:"type:date;not null"`
	Gender        Genders        `gorm:"type:enum('male','female');not null"`
	PlanTypeID    uint           `gorm:"not null"`
	Dependence    *string        // VARCHAR bisa *string jika NULLABLE, atau string jika NOT NULL
	BankNumber    string         `gorm:"not null"`
	ProRate       float64        `gorm:"not null"`
	JoinDate      time.Time      `gorm:"type:date;not null"`
	Patient       Patient        `gorm:"foreignKey:EmployeeID"`
	Department    Department     `gorm:"foreignKey:DepartmentID"`
	PlanType      PlanType       `gorm:"foreignKey:PlanTypeID"`
	FamilyMembers []FamilyMember `gorm:"foreignKey:EmployeeID"`
	Claims        []Claim        `gorm:"foreignKey:EmployeeID"`
}

type FamilyMember struct {
	ID               uint              `gorm:"primaryKey;autoIncrement"`
	EmployeeID       uint              `gorm:"not null"`
	Name             string            `gorm:"not null"`
	RelationshipType RelationshipTypes `gorm:"type:enum('father', 'mother', 'wife', 'husband', 'child');not null"`
	PlanTypeID       uint              `gorm:"not null"`
	BirthDate        time.Time         `gorm:"type:date;not null"`
	Gender           Genders           `gorm:"type:enum('male','female');not null"`
	Patient          Patient           `gorm:"foreignKey:FamilyMemberID"`
	Employee         *Employee         `gorm:"foreignKey:EmployeeID"`
	PlanType         PlanType          `gorm:"foreignKey:PlanTypeID"`
}

type Patient struct {
	ID              uint             `gorm:"primaryKey;autoIncrement"`
	Name            string           `gorm:"not null"`
	BirthDate       time.Time        `gorm:"type:date;not null"`
	Gender          Genders          `gorm:"type:enum('male','female');not null"`
	EmployeeID      *uint            `gorm:"uniqueIndex"`
	FamilyMemberID  *uint            `gorm:"uniqueIndex"`
	PlanTypeID      uint             `gorm:"not null"`
	Employee        *Employee        `gorm:"foreignKey:EmployeeID"`
	FamilyMember    *FamilyMember    `gorm:"foreignKey:FamilyMemberID"`
	Claims          []Claim          `gorm:"foreignKey:PatientID"`
	PatientBenefits []PatientBenefit `gorm:"foreignKey:PatientID"`
	PlanType        PlanType         `gorm:"foreignKey:PlanTypeID"`
	Benefits        []Benefit        `gorm:"many2many:patient_benefits;joinForeignKey:PatientID;joinReferences:BenefitID"`
}

type YearlyBenefitClaim struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	Code        string     `gorm:"column:code;not null"`
	YearlyClaim float64    `gorm:"column:yearly_claim;not null"`
	CreatedAt   time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoCreateTime;autoUpdateTime"`
	Benefits    []*Benefit `gorm:"foreignKey:YearlyBenefitClaimID"`
}

func (e YearlyBenefitClaim) TableName() string {
	return "yearly_benefit_claims"
}

type Benefit struct {
	ID                   uint                `gorm:"primaryKey;autoIncrement"`
	Name                 string              `gorm:"not null"`
	PlanTypeID           uint                `gorm:"not null"`
	YearlyBenefitClaimID *uint               `gorm:"column:yearly_benefit_claim_id"`
	Detail               *string             `gorm:"column:detail"`
	Code                 string              `gorm:"unique;not null"`
	LimitationType       LimitationType      `gorm:"enum('Per Day', 'Per Month', 'Per Year', 'Per Incident', 'Per Pregnancy');not null"`
	Plafond              *float64            `gorm:"column:plafond"`
	YearlyMax            *float64            `gorm:"column:yearly_max"`
	PlanType             PlanType            `gorm:"foreignKey:PlanTypeID"`
	PatientBenefits      []PatientBenefit    `gorm:"foreignKey:BenefitID"` // Ini sudah benar
	YearlyBenefitClaim   *YearlyBenefitClaim `gorm:"foreignKey:YearlyBenefitClaimID"`
	Patients             []Patient           `gorm:"many2many:patient_benefits;joinForeignKey:BenefitID;joinReferences:PatientID"`
}

type PatientBenefit struct {
	ID               uint                 `gorm:"primaryKey;autoIncrement"`
	PatientID        uint                 `gorm:"not null"`
	BenefitID        uint                 `gorm:"not null"`
	RemainingPlafond *float64             `gorm:"type:decimal(10,2);"`
	InitialPlafond   *float64             `gorm:"type:decimal(10,2);"`
	YearlyMax        *float64             `gorm:"column:yearly_max"`
	StartDate        time.Time            `gorm:"type:date;not null"`
	EndDate          *time.Time           `gorm:"type:date"`
	Status           PatientBenefitStatus `gorm:"type:enum('active','exhausted','expired');default:'active'"`
	CreatedAt        time.Time            `gorm:"not null;autoCreateTime"`
	UpdatedAt        *time.Time           `gorm:"autoUpdateTime"`

	Patient Patient `gorm:"foreignKey:PatientID"`
	Benefit Benefit `gorm:"foreignKey:BenefitID"`
	Claims  []Claim `gorm:"foreignKey:PatientBenefitID"`
}

type Claim struct {
	ID                  uint        `gorm:"primaryKey;autoIncrement"`
	PatientBenefitID    uint        `gorm:"not null"`
	PatientID           uint        `gorm:"not null"`
	EmployeeID          uint        `gorm:"not null"`
	ClaimAmount         float64     `gorm:"type:decimal(10,2);not null"`
	TransactionTypeID   *uint       `gorm:"null"`
	TransactionDate     *time.Time  `gorm:"type:date;null"`
	SubmissionDate      *time.Time  `gorm:"type:date;null"`
	SLA                 *SLA        `gorm:"type:enum('meet','overdue');null"`
	ApprovedAmount      *float64    `gorm:"type:decimal(10,2);null"`
	ClaimStatus         ClaimStatus `gorm:"type:enum('On Plafond','Over Plafond');not null"`
	MedicalFacilityName *string
	City                *string
	Diagnosis           *string
	DocLink             *string
	TransactionStatus   TransactionStatus `gorm:"type:enum('Successful','Pending','Failed');not null"`
	CreatedAt           time.Time         `gorm:"not null;autoCreateTime"`
	UpdatedAt           *time.Time        `gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt           *gorm.DeletedAt   `gorm:"index"`

	Patient         Patient          `gorm:"foreignKey:PatientID"`
	Employee        Employee         `gorm:"foreignKey:EmployeeID"`
	PatientBenefit  PatientBenefit   `gorm:"foreignKey:PatientBenefitID"`
	TransactionType *TransactionType `gorm:"foreignKey:TransactionTypeID"`
}
