package seed

import (
	"errors"
	"log"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"gorm.io/gorm"
)

func SeedTransactionTypes(db *gorm.DB) {
	transactionTypes := []entity.TransactionType{
		{Name: "Invoice"},
		{Name: "Reimbursement"},
		{Name: "Advance"},
		{Name: "Credit Card"},
	}

	for _, trx := range transactionTypes {
		var existingTrx entity.TransactionType
		if err := db.Where("name = ?", trx.Name).First(&existingTrx).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&trx).Error; err != nil {
					log.Printf("Error seeding transaction type %s: %v\n", trx.Name, err)
				} else {
					log.Printf("Transaction type %s seeded successfully.\n", trx.Name)
				}
			} else {
				log.Printf("Error checking transaction type %s: %v\n", trx.Name, err)
			}
		} else {
			log.Printf("Transaction type %s already exists, skipping.\n", existingTrx.Name)
		}
	}
}

func SeedFamilyMemberAndEmployee(db *gorm.DB) {
	var department entity.Department
	if db.Where("name = ?", "IT Support").First(&department).Error != nil {
		log.Fatalf("Department IT Support not found")
	}
	var planType entity.PlanType
	if db.Where("name = ?", "PLAN A").First(&planType).Error != nil {
		log.Fatalf("PlanType 'PLAN A' not found")
	}
	birthDate, errParse := time.Parse("2006-01-02", "2006-08-17")
	if errParse != nil {
		log.Fatalf("Error parsing birth date: %v", errParse)
	}
	dependence := "One Child, One Wife"
	employees := []entity.Employee{
		{
			Name:         "John Doe",
			DepartmentID: department.ID,
			Position:     department.Name,
			Email:        "johndoe@gmail.com",
			Phone:        "+1 2887 2982 2394",
			BirthDate:    birthDate,
			Gender:       "male",
			PlanTypeID:   planType.ID,
			Dependence:   &dependence,
			BankNumber:   "1234-456-789",
			ProRate:      70.00,
		},
		{
			Name:         "Mary Jane",
			DepartmentID: department.ID,
			Position:     department.Name,
			Email:        "maryjane@gmail.com",
			Phone:        "+1 2878 2982 2394",
			BirthDate:    birthDate,
			Gender:       "female",
			PlanTypeID:   planType.ID,
			Dependence:   &dependence,
			BankNumber:   "1243-456-789",
			ProRate:      70.00,
		},
	}

	for _, employee := range employees {
		var existingEmpl entity.Employee
		if err := db.Where("name = ?", employee.Name).First(&existingEmpl).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&employee).Error; err != nil {
					log.Printf("Error seeding employee %s: %v\n", employee.Name, err)
				} else {
					if db.Create(&entity.FamilyMember{
						EmployeeID: employee.ID,
						Name:       employee.Name,
						RelationshipType: func(gender entity.Genders) entity.RelationshipTypes {
							if gender == entity.GenderMale {
								return "husband"
							} else {
								return "wife"
							}
						}(employee.Gender),
						PlanTypeID: employee.PlanTypeID,
						BirthDate:  employee.BirthDate,
						Gender:     employee.Gender,
					}).Error != nil {
						log.Printf("Error seeding employee %s: %v\n", employee.Name, err)
					}
					log.Printf("Employee %s seeded successfully.\n", employee.Name)
				}
			} else {
				log.Printf("Error checking employee %s: %v\n", employee.Name, err)
			}
		} else {
			log.Printf("Employee %s already exists, skipping.\n", existingEmpl.Name)
		}
	}
}

func SeedClaimsAndPatients(db *gorm.DB) {
	for i := 0; i < 100; i++ {
		genders := []string{"male", "female"}
		relationshipTypes := []string{"husband", "wife", "child", "father", "mother"}
		benefits := make([]entity.Benefit, 10)
		if errFindBenefits := db.Find(&benefits).Limit(10).Error; errFindBenefits != nil {
			log.Fatalf("Error find benefits: %v\n", errFindBenefits)
		}
		planTypes := make([]entity.PlanType, 4)
		if errFindPlanTypes := db.Find(&planTypes).Limit(4).Error; errFindPlanTypes != nil {
			log.Fatalf("Error find planTypes: %v\n", errFindPlanTypes)
		}

		transactionTypes := make([]entity.TransactionType, 4)
		if errFindTransactionTypes := db.Find(&transactionTypes).Limit(4).Error; errFindTransactionTypes != nil {
			log.Fatalf("Error find transactionTypes: %v\n", errFindTransactionTypes)
		}

		departments := make([]entity.Department, 4)
		if errFindDepartment := db.Find(&departments).Limit(4).Error; errFindDepartment != nil {
			log.Fatalf("Error find departments: %v\n", errFindDepartment)
		}

		randomBenefit := helper.RandomInt(0, len(benefits)-1)
		randomDepartment := helper.RandomInt(0, len(departments)-1)
		randomPlanType := helper.RandomInt(0, len(planTypes)-1)
		randomTransactionType := helper.RandomInt(0, len(transactionTypes)-1)

		fakeDate, errParseDate := time.Parse("2006-01-02", faker.Date())
		if errParseDate != nil {
			log.Fatalf("Error seeding fakeDate: %v\n", errParseDate)
		}
		employee := entity.Employee{
			Name:         faker.Name(),
			DepartmentID: departments[randomDepartment].ID,
			Position:     departments[randomDepartment].Name,
			Email:        faker.Email(),
			Phone:        faker.Phonenumber(),
			BirthDate:    fakeDate,
			Gender:       entity.Genders(genders[helper.RandomInt(0, len(genders)-1)]),
			PlanTypeID:   planTypes[randomPlanType].ID,
			Dependence:   ptrString(faker.Word()),
			BankNumber:   faker.CreditCardNumber,
			ProRate:      helper.ProRateRemainingMonthsPercent(time.Now(), time.Now()),
			JoinDate:     time.Now(),
		}
		if err := db.Create(&employee).Error; err != nil {
			log.Fatalf("Error when seeding employee: %v\n", err)
		}
		patient := entity.Patient{
			Name:       faker.Name(),
			BirthDate:  fakeDate,
			Gender:     entity.Genders(genders[helper.RandomInt(0, len(genders)-1)]),
			PlanTypeID: planTypes[randomPlanType].ID,
			EmployeeID: &employee.ID,
			FamilyMember: &entity.FamilyMember{
				Name:             faker.Name(),
				RelationshipType: entity.RelationshipTypes(relationshipTypes[helper.RandomInt(0, len(relationshipTypes)-1)]),
				PlanTypeID:       planTypes[randomPlanType].ID,
				BirthDate:        fakeDate,
				Gender:           entity.Genders(genders[helper.RandomInt(0, len(genders)-1)]),
				EmployeeID:       employee.ID,
			},
		}
		if err := db.Create(&patient).Error; err != nil {
			log.Fatalf("Error when seeding patient: %v\n", err)
		}
		patientBenefit := entity.PatientBenefit{
			PatientID:        patient.ID,
			BenefitID:        benefits[randomBenefit].ID,
			RemainingPlafond: benefits[randomBenefit].Plafond,
			InitialPlafond:   benefits[randomBenefit].Plafond,
			StartDate:        employee.JoinDate,
			EndDate:          &employee.JoinDate,
			Status:           entity.PatientBenefitStatusActive,
		}
		if err := db.Create(&patientBenefit).Error; err != nil {
			log.Fatalf("Error when seeding patient: %v\n", err)
		}
		sla := []entity.SLA{entity.SLAOverdue, entity.SLAMeet}[helper.RandomInt(0, 1)]
		claimStatus := []entity.TransactionStatus{entity.TransactionStatusSuccessful, entity.TransactionStatusFailed, entity.TransactionStatusPending}[helper.RandomInt(0, 2)]
		db.Create(&entity.Claim{
			PatientBenefitID: patientBenefit.ID,
			PatientID:        patient.ID,
			EmployeeID:       employee.ID,
			ClaimAmount: func() float64 {
				if benefits[randomBenefit].Plafond != nil {
					return *benefits[randomBenefit].Plafond
				}
				return float64(helper.RandomInt(1000000, 10000000))
			}(),
			TransactionTypeID: &transactionTypes[randomTransactionType].ID,
			TransactionDate:   ptrDate(time.Now()),
			SubmissionDate:    ptrDate(time.Now()),
			SLA:               &sla,
			ApprovedAmount: func() *float64 {
				if benefits[randomBenefit].Plafond != nil {
					return benefits[randomBenefit].Plafond
				}
				return ptrFloat64(float64(helper.RandomInt(1000000, 10000000)))
			}(),
			ClaimStatus:         entity.ClaimStatusOnPlafond,
			MedicalFacilityName: ptrString(faker.ChineseName()),
			City:                ptrString("Daerah Istimewa Yogyakarta"),
			Diagnosis:           ptrString(faker.Word()),
			DocLink:             ptrString(faker.URL()),
			TransactionStatus:   claimStatus,
		})
	}
}
