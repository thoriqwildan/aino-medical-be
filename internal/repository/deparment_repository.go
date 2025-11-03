package repository

import (
	"errors"

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

func (dr *DepartmentRepository) FindOrCreate(db *gorm.DB, department *entity.Department) (*entity.Department, error) {
	var result entity.Department
	errTake := db.Model(entity.Department{}).Where("name = ?", department.Name).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		dr.Log.Error("Error when find department in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil 
	}

	errCreate := db.Create(department).Error
	if errCreate != nil {
		dr.Log.Error("Error when create department in method find or create")
		return nil, errCreate
	}

	return department, nil
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
	if request.Limit > 0 {
		baseQuery = baseQuery.Limit(request.Limit)
	}
	if request.Page > 0 {
		baseQuery = baseQuery.Offset((request.Page - 1) * request.Limit)
	}
	err := baseQuery.
		Find(&departments).Error
	if err != nil {
		return nil, 0, err
	}

	return departments, total, nil
}
