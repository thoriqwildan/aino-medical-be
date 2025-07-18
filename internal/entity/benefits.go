package entity

import "time"

type PlanType struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"unique;not null"`
	Description *string
	Benefits    []Benefit    `gorm:"foreignKey:PlanTypeID"`
	Employees   []Employee   `gorm:"foreignKey:PlanTypeID"`
	FamilyMembers []FamilyMember `gorm:"foreignKey:PlanTypeID"`
}

// TransactionType represents the transaction_types table.
type TransactionType struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"unique;not null"`
	Claims []Claim `gorm:"foreignKey:TransactionTypeID"`
}

// LimitationType represents the limitation_types table.
type LimitationType struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"unique;not null"`
	Benefits []Benefit `gorm:"foreignKey:LimitationTypeID"`
}

// Department represents the departments table.
type Department struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"` // For soft delete
	Employees []Employee `gorm:"foreignKey:DepartmentID"`
}

// Patient represents the patients table.
type Patient struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"not null"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	Gender      Genders   `gorm:"type:enum('male','female');not null"`
	Employee    *Employee `gorm:"foreignKey:PatientID"` // A patient might be an employee
	FamilyMember *FamilyMember `gorm:"foreignKey:PatientID"` // A patient might be a family member
	Claims      []Claim   `gorm:"foreignKey:PatientID"`
}

// Employee represents the employees table.
type Employee struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	Name          string    `gorm:"not null"`
	PatientID     uint      `gorm:"unique;not null"`
	DepartmentID  uint      `gorm:"not null"`
	Position      string    `gorm:"not null"`
	Email         string    `gorm:"unique;not null"`
	Phone         string    `gorm:"not null"`
	BirthDate     time.Time `gorm:"type:date;not null"`
	Gender        Genders   `gorm:"type:enum('male','female');not null"`
	PlanTypeID    uint      `gorm:"not null"`
	Dependence    *string
	BankNumber    string    `gorm:"not null"`
	JoinDate      time.Time `gorm:"type:date;not null"`
	Patient       Patient   `gorm:"foreignKey:PatientID"` // Belongs To Patient
	Department    Department `gorm:"foreignKey:DepartmentID"` // Belongs To Department
	PlanType      PlanType  `gorm:"foreignKey:PlanTypeID"` // Belongs To PlanType
	FamilyMembers []FamilyMember `gorm:"foreignKey:EmployeeID"` // Has Many FamilyMembers
	Claims        []Claim   `gorm:"foreignKey:EmployeeID"` // Has Many Claims
}

// FamilyMember represents the family_members table.
type FamilyMember struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	PatientID   uint      `gorm:"unique;not null"`
	EmployeeID  uint      `gorm:"not null"`
	Name        string    `gorm:"not null"`
	PlanTypeID  uint      `gorm:"not null"`
	BirthDate   time.Time `gorm:"type:date;not null"`
	Gender      Genders   `gorm:"type:enum('male','female');not null"`
	Patient     Patient   `gorm:"foreignKey:PatientID"` // Belongs To Patient
	Employee    Employee  `gorm:"foreignKey:EmployeeID"` // Belongs To Employee
	PlanType    PlanType  `gorm:"foreignKey:PlanTypeID"` // Belongs To PlanType
}

// Benefit represents the benefits table.
type Benefit struct {
	ID               uint           `gorm:"primaryKey;autoIncrement"`
	Name             string         `gorm:"not null"`
	PlanTypeID       uint           `gorm:"not null"`
	Detail           *string
	Code             string         `gorm:"unique;not null"`
	LimitationTypeID uint           `gorm:"not null"`
	Plafond          int            `gorm:"not null"`
	YearlyMax        int            `gorm:"not null"`
	PlanType         PlanType       `gorm:"foreignKey:PlanTypeID"` // Belongs To PlanType
	LimitationType   LimitationType `gorm:"foreignKey:LimitationTypeID"` // Belongs To LimitationType
	Claims           []Claim        `gorm:"foreignKey:BenefitID"`
}

// Claim represents the claims table.
type Claim struct {
	ID                  uint            `gorm:"primaryKey;autoIncrement"`
	PatientID           uint            `gorm:"not null"`
	EmployeeID          uint            `gorm:"not null"`
	BenefitID           uint            `gorm:"not null"`
	ClaimAmount         float64         `gorm:"type:decimal(10,2);not null"`
	TransactionTypeID   uint            `gorm:"not null"`
	TransactionDate     time.Time       `gorm:"type:date;not null"`
	SubmissionDate      time.Time       `gorm:"type:date;not null"`
	SLA                 SLA             `gorm:"type:enum('meet','overdue');not null"`
	ApprovedAmount      float64         `gorm:"type:decimal(10,2);not null"`
	ClaimStatus         ClaimStatus     `gorm:"type:enum('On Plafond','Over Plafond');not null"`
	MedicalFacilityName *string
	City                *string
	Diagnosis           *string
	DocLink             *string
	TransactionStatus   TransactionStatus `gorm:"type:enum('Successful','Pending','Failed');not null"`
	CreatedAt           time.Time       `gorm:"not null;autoCreateTime"`
	UpdatedAt           *time.Time       `gorm:"autoUpdateTime"`
	DeletedAt           *time.Time       `gorm:"index"` // For soft delete

	Patient         Patient         `gorm:"foreignKey:PatientID"` // Belongs To Patient
	Employee        Employee        `gorm:"foreignKey:EmployeeID"` // Belongs To Employee
	Benefit         Benefit         `gorm:"foreignKey:BenefitID"` // Belongs To Benefit
	TransactionType TransactionType `gorm:"foreignKey:TransactionTypeID"` // Belongs To TransactionType
}