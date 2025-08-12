package seed

import (
	"log"

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