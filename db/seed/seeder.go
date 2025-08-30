package seed

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) {
	log.Println("Running all seeders...")

	SeedTransactionTypes(db)
	SeedDepartments(db)
	SeedBenefits(db)
	SeedFamilyMemberAndEmployee(db)
	SeedClaimsAndPatients(db)
	log.Println("Database seeding completed successfully.")
}
