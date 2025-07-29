package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

type EmployeeRepository struct {
	Repository[entity.Employee]
	Log *logrus.Logger
}

func NewEmployeeRepository(log *logrus.Logger) *EmployeeRepository {
	return &EmployeeRepository{
		Log: log,
	}
}

func (er *EmployeeRepository) GetByEmail(db *gorm.DB, email string) error {
	return db.Where("email = ?", email).First(&entity.Employee{}).Error
}