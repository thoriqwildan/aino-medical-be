package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	Repository[entity.Department]
	Log *logrus.Logger
}

func NewDepartmentRepository(log *logrus.Logger) *DepartmentRepository {
	return &DepartmentRepository{
		Log: log,
	}
}

func (dr *DepartmentRepository) GetByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.Department{}).Error
}