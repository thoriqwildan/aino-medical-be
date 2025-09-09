package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
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

func (er *EmployeeRepository) GetDepartmentByID(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).First(&entity.Department{}).Error
}

func (er *EmployeeRepository) FindById(db *gorm.DB, id uint, employee *entity.Employee) error {
	return db.Where("id = ?", id).
		Preload("Department").
		Preload("PlanType").
		Preload("FamilyMembers").
		Preload("FamilyMembers.PlanType").
		First(employee).Error
}

func (er *EmployeeRepository) SearchEmployees(db *gorm.DB, request *model.SearchPagingQuery) ([]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	baseQuery := db.Model(&entity.Employee{})

	if request.SearchValue != "" {
		baseQuery.Where("name LIKE ?", "%"+request.SearchValue+"%")
	}

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
		Preload("Department").
		Preload("PlanType").
		Preload("FamilyMembers").
		Preload("FamilyMembers.PlanType").
		Order("name ASC").
		Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}
