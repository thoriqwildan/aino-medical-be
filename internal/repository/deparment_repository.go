package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
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

func (r *DepartmentRepository) SearchDepartments(db *gorm.DB, request *model.PagingQuery) ([]entity.Department, int64, error) {
	var departments []entity.Department
	var total int64

	baseQuery := db.Model(&entity.Department{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Find(&departments).Error
	if err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}