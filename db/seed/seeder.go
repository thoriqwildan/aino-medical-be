package seed

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) {
	log.Println("Running all seeders...")

	SeedLimitationTypes(db)
	SeedPlanTypes(db)
	SeedTransactionTypes(db)
	SeedDepartments(db)

	log.Println("Database seeding completed successfully.")
}