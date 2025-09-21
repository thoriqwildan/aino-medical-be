package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type FamilyMemberRepository struct {
	Repository[entity.FamilyMember]
	Log *logrus.Logger
}

func NewFamilyMemberRepository(log *logrus.Logger) *FamilyMemberRepository {
	return &FamilyMemberRepository{
		Log: log,
	}
}

func (r *FamilyMemberRepository) GetByID(db *gorm.DB, familyMember *entity.FamilyMember, id any) error {
	return db.Where("id = ?", id).
		Preload("PlanType").
		Preload("Employee").
		Preload("Employee.PlanType").
		Preload("Employee.Department").
		First(familyMember).Error
}

func (r *FamilyMemberRepository) GetByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.FamilyMember{}).Error
}

func (r *FamilyMemberRepository) GetEmployeeById(db *gorm.DB, employee *entity.Employee, id any) error {
	return db.Where("id = ?", id).First(employee).Error
}

func (r *FamilyMemberRepository) SearchFamilyMember(db *gorm.DB, request *model.SearchPagingQuery) ([]entity.FamilyMember, int64, error) {
	var familyMembers []entity.FamilyMember
	var total int64

	baseQuery := db.Model(&entity.FamilyMember{})

	if request.SearchValue != "" {
		baseQuery.Where("name LIKE ?", "%"+request.SearchValue+"%")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Preload("PlanType").
		Preload("Employee").
		Preload("Employee.PlanType").
		Preload("Employee.Department").
		Order("name ASC").
		Find(&familyMembers).Error
	if err != nil {
		return nil, 0, err
	}

	return familyMembers, total, nil
}
