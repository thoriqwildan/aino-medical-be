package seed

import (
	"errors"
	"log"

	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

func SeedDepartments(db *gorm.DB) {
	departments := []entity.Department{
		{Name: "Human Resource"},
		{Name: "Finance"},
		{Name: "IT Support"},
		{Name: "Human Capital"},
	}

	for _, dept := range departments {
		var existingDept entity.Department
		if err := db.Where("name = ?", dept.Name).First(&existingDept).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&dept).Error; err != nil {
					log.Printf("Error seeding department %s: %v\n", dept.Name, err)
				} else {
					log.Printf("department %s seeded successfully.\n", dept.Name)
				}
			} else {
				log.Printf("Error checking department %s: %v\n", dept.Name, err)
			}
		} else {
			log.Printf("department %s already exists, skipping.\n", existingDept.Name)
		}
	}
}
