package seed

import (
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

func SeedClaimsAndPatients(db *gorm.DB) {

	for i := 0; i < 100; i++ {

		genders := []string{"male", "female"}
		relationshipTypes := []string{"husband", "wife", "child", "father", "mother"}
		var benefits []entity.Benefit
		if err := db.Preload("YearlyBenefitClaim").Limit(10).Find(&benefits).Error; err != nil {
			log.Fatalf("Error find benefits: %v\n", err)
		}
		if len(benefits) == 0 {
			log.Fatal("No benefits found")
		}

		var planTypes []entity.PlanType
		if err := db.Limit(4).Find(&planTypes).Error; err != nil {
			log.Fatalf("Error find planTypes: %v\n", err)
		}
		if len(planTypes) == 0 {
			log.Fatal("No planTypes found")
		}

		var transactionTypes []entity.TransactionType
		if err := db.Limit(4).Find(&transactionTypes).Error; err != nil {
			log.Fatalf("Error find transactionTypes: %v\n", err)
		}

		var departments []entity.Department
		if err := db.Limit(4).Find(&departments).Error; err != nil {
			log.Fatalf("Error find departments: %v\n", err)
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
			ProRate:      helper.ProRateRemainingMonthsPercent(time.Now(), time.Now().AddDate(0, 2, 0)),
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
		if benefits[randomBenefit].YearlyBenefitClaimID != nil && benefits[randomBenefit].YearlyBenefitClaim != nil {
			var count int64
			errFind := db.Model(entity.PatientYearlyBenefitClaim{}).
				Where("yearly_benefit_claim_id = ?", *benefits[randomBenefit].YearlyBenefitClaimID).
				Where("patient_id = ?", patientBenefit.PatientID).
				Count(&count).Error
			if errFind != nil {
				log.Fatalf("Error finding patient yearly benefit claims: %v\n", errFind)
			}
			if count == 0 {
				if errCreate := db.Model(entity.PatientYearlyBenefitClaim{}).
					Create(&entity.PatientYearlyBenefitClaim{
						PatientID:            patientBenefit.PatientID,
						YearlyBenefitClaimID: *benefits[randomBenefit].YearlyBenefitClaimID,
						YearlyClaimRemaining: benefits[randomBenefit].YearlyBenefitClaim.YearlyClaim,
					}).Error; errCreate != nil {
					log.Fatalf("Error creating patient yearly benefit claims: %v\n", errCreate)
				}
			}
		}
		if err := db.Create(&patientBenefit).Error; err != nil {
			log.Fatalf("Error when seeding patient: %v\n", err)
		}
		sla := []entity.SLA{entity.SLAOverdue, entity.SLAMeet}[helper.RandomInt(0, 1)]
		claimStatus := []entity.TransactionStatus{entity.TransactionStatusSuccessful, entity.TransactionStatusFailed, entity.TransactionStatusPending}[helper.RandomInt(0, 2)]
		startRandom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		endRandom := time.Now()
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
			TransactionDate:   ptrDate(helper.RandomDate(startRandom, endRandom)),
			SubmissionDate:    ptrDate(helper.RandomDate(startRandom, endRandom)),
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
