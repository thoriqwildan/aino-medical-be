package seed

import (
	"errors"
	"log"
	"time"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
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

func SeedLimitationTypes(db *gorm.DB) {
	limitationTypes := []entity.LimitationType{
		{Name: "Annual"},
		{Name: "Per Incident"},
		{Name: "Per Pregnancy"},
	}

	for _, lt := range limitationTypes {
		var existingLt entity.LimitationType
		if err := db.Where("name = ?", lt.Name).First(&existingLt).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&lt).Error; err != nil {
					log.Printf("Error seeding limitation type %s: %v\n", lt.Name, err)
				} else {
					log.Printf("Limitation type %s seeded successfully.\n", lt.Name)
				}
			} else {
				log.Printf("Error checking limitation type %s: %v\n", lt.Name, err)
			}
		} else {
			log.Printf("Limitation type %s already exists, skipping.\n", existingLt.Name)
		}
	}
}

func SeedPlanTypes(db *gorm.DB) {
	planTypes := []entity.PlanType{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
	}

	for _, pt := range planTypes {
		var existingPt entity.PlanType
		if err := db.Where("name = ?", pt.Name).First(&existingPt).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pt).Error; err != nil {
					log.Printf("Error seeding plan type %s: %v\n", pt.Name, err)
				} else {
					log.Printf("Plan type %s seeded successfully.\n", pt.Name)
				}
			} else {
				log.Printf("Error checking plan type %s: %v\n", pt.Name, err)
			}
		} else {
			log.Printf("Plan type %s already exists, skipping.\n", existingPt.Name)
		}
	}
}

func SeedFamilyMemberAndEmployee(db *gorm.DB) {
	var department entity.Department
	if db.Where("name = ?", "IT Support").First(&department).Error != nil {
		log.Fatalf("Department IT Support not found")
	}
	var planType entity.PlanType
	if db.Where("name = ?", "A").First(&planType).Error != nil {
		log.Fatalf("PlanType 'A' not found")
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
