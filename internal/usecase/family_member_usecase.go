package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/model/converter"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type FamilyMemberUseCase struct {
	Repository *repository.FamilyMemberRepository
	DB *gorm.DB
	Validate *validator.Validate
	Log *logrus.Logger
}

func NewFamilyMemberUseCase(repo *repository.FamilyMemberRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger) *FamilyMemberUseCase {
	return &FamilyMemberUseCase{
		Repository: repo,
		DB: db,
		Validate: validate,
		Log: log,
	}
}

func (uc *FamilyMemberUseCase) Create(ctx context.Context, request *model.FamilyMemberRequest) (*model.FamilyMemberResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation failed for FamilyMemberRequest")
		return nil, err
	}

	employee := &entity.Employee{}
	if err := uc.Repository.GetEmployeeById(tx, employee, request.EmployeeID); err != nil {
		uc.Log.WithError(err).Error("Employee not found")
		return nil, err
	}

	familyMember := &entity.FamilyMember{
		EmployeeID: request.EmployeeID,
		Name: 		 request.Name,
		PlanTypeID: employee.PlanTypeID,
		BirthDate: time.Time(request.BirthDate),
		Gender: entity.Genders(request.Gender),
		Patient: entity.Patient{
			Name: request.Name,
			BirthDate: time.Time(request.BirthDate),
			Gender: entity.Genders(request.Gender),
		},
	}

	if err := tx.Create(familyMember).Error; err != nil {
		uc.Log.WithError(err).Error("Failed to create family member")
		return nil, err
	}

	if err := uc.Repository.GetByID(tx, familyMember, familyMember.ID); err != nil {
		uc.Log.WithError(err).Error("Failed to retrieve created family member")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}

	return converter.FamilyMemberToResponse(familyMember), nil
}

func (uc *FamilyMemberUseCase) GetByID(ctx context.Context, id uint) (*model.FamilyMemberResponse, error) {
	tx := uc.DB.WithContext(ctx)

	familyMember := &entity.FamilyMember{}
	if err := uc.Repository.GetByID(tx, familyMember, id); err != nil {
		uc.Log.WithError(err).Error("Failed to get family member by ID")
		return nil, fiber.NewError(fiber.StatusNotFound, "Family member not found")
	}

	return converter.FamilyMemberToResponse(familyMember), nil
}

func (uc *FamilyMemberUseCase) GetAll(ctx context.Context, request *model.PagingQuery) ([]model.FamilyMemberResponse, int64, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation error in GetAllFamilyMembers")
		return nil, 0, err
	}

	familyMembers, total, err := uc.Repository.SearchFamilyMember(tx, request)
	if err != nil {
		uc.Log.WithError(err).Error("Error searching family members")
		return nil, 0, err
	}

	responses := make([]model.FamilyMemberResponse, len(familyMembers))
	for i, fm := range familyMembers {
		responses[i] = *converter.FamilyMemberToResponse(&fm)
	}
	return responses, total, nil
}

func (uc *FamilyMemberUseCase) Update(ctx context.Context, request *model.UpdateFamilyMemberRequest) (*model.FamilyMemberResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.WithError(err).Error("Validation failed for UpdateFamilyMemberRequest")
		return nil, err
	}

	familyMember := &entity.FamilyMember{}
	if err := uc.Repository.GetByID(tx, familyMember, request.ID); err != nil {
		uc.Log.WithError(err).Error("Family member not found")
		return nil, fiber.NewError(fiber.StatusNotFound, "Family member not found")
	}


	familyMember.Name = request.Name
	familyMember.BirthDate = time.Time(request.BirthDate)
	familyMember.Gender = entity.Genders(request.Gender)
	familyMember.Patient = entity.Patient{
		Name: request.Name,
		BirthDate: time.Time(request.BirthDate),
		Gender: entity.Genders(request.Gender),
	}
	if err := uc.Repository.Update(tx, familyMember); err != nil {
		uc.Log.WithError(err).Error("Failed to update family member")
		return nil, err
	}
	if err := uc.Repository.GetByID(tx, familyMember, familyMember.ID); err != nil {
		uc.Log.WithError(err).Error("Failed to retrieve updated family member")
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return nil, err
	}
	return converter.FamilyMemberToResponse(familyMember), nil
}

func (uc *FamilyMemberUseCase) Delete(ctx context.Context, id uint) error {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	familyMember := &entity.FamilyMember{}
	if err := uc.Repository.GetByID(tx, familyMember, id); err != nil {
		uc.Log.WithError(err).Error("Family member not found")
		return fiber.NewError(fiber.StatusNotFound, "Family member not found")
	}

	if err := tx.Delete(familyMember).Error; err != nil {
		uc.Log.WithError(err).Error("Failed to delete family member")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.WithError(err).Error("Failed to commit transaction")
		return err
	}

	return nil
}
