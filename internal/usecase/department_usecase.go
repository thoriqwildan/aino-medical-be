package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/model/converter"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type DepartmentUseCase struct {
	DepartmentRepository *repository.DepartmentRepository
	Validate *validator.Validate
	Log *logrus.Logger
	DB *gorm.DB
}

func NewDepartmentUseCase(repo *repository.DepartmentRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) *DepartmentUseCase {
	return &DepartmentUseCase{
		DepartmentRepository: repo,
		DB:                   db,
		Log:                  log,
		Validate:             validate,
	}
}

func (du *DepartmentUseCase) Create(ctx context.Context, request *model.DepartmentRequest) (*model.DepartmentResponse, error) {
	tx := du.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := du.Validate.Struct(request); err != nil {
		du.Log.WithError(err).Error("Validation error in CreateDepartment")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	if err := du.DepartmentRepository.GetByName(tx, request.Name); err == nil {
		du.Log.WithField("name", request.Name).Error("Department with this name already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Department with this name already exists")
	}

	department := &entity.Department{
		Name: request.Name,
	}

	if err := du.DepartmentRepository.Create(tx, department); err != nil {
		du.Log.WithError(err).Error("Error creating department")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create department")
	}

	if err := tx.Commit().Error; err != nil {
		du.Log.WithError(err).Error("Error committing transaction in CreateDepartment")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return converter.DepartmentToResponse(department), nil
}

func (du *DepartmentUseCase) GetById(ctx context.Context, id uint) (*model.DepartmentResponse, error) {
	tx := du.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	department := &entity.Department{}
	if err := du.DepartmentRepository.FindById(tx, department, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			du.Log.WithField("id", id).Error("Department not found")
			return nil, fiber.NewError(fiber.StatusNotFound, "Department not found")
		}
		du.Log.WithError(err).Error("Error retrieving department by ID")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve department")
	}

	if err := tx.Commit().Error; err != nil {
		du.Log.WithError(err).Error("Error committing transaction in GetById")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return converter.DepartmentToResponse(department), nil
}

func (du *DepartmentUseCase) GetAll(ctx context.Context, request *model.PagingQuery) ([]model.DepartmentResponse, int64, error) {
	tx := du.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := du.Validate.Struct(request); err != nil {
		du.Log.WithError(err).Error("Validation error in GetAllDepartments")
		return nil, 0, err
	}

	departments, total, err := du.DepartmentRepository.SearchDepartments(tx, request)
	if err != nil {
		du.Log.WithError(err).Error("Error searching departments")
		return nil, 0, err
	}

	responses := make([]model.DepartmentResponse, len(departments))
	for i, d := range departments {
		responses[i] = *converter.DepartmentToResponse(&d)
	}
	return responses, total, nil
}

func (du *DepartmentUseCase) Update(ctx context.Context, request *model.UpdateDepartmentRequest) (*model.DepartmentResponse, error) {
	tx := du.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := du.Validate.Struct(request); err != nil {
		du.Log.WithError(err).Error("Validation error in UpdateDepartment")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	department := &entity.Department{}
	if err := du.DepartmentRepository.FindById(tx, department, request.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			du.Log.WithField("id", request.ID).Error("Department not found")
			return nil, fiber.NewError(fiber.StatusNotFound, "Department not found")
		}
		du.Log.WithError(err).Error("Error retrieving department for update")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve department for update")
	}

	department.Name = request.Name

	if err := du.DepartmentRepository.Update(tx, department); err != nil {
		du.Log.WithError(err).Error("Error updating department")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update department")
	}

	if err := tx.Commit().Error; err != nil {
		du.Log.WithError(err).Error("Error committing transaction in UpdateDepartment")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	return converter.DepartmentToResponse(department), nil
}

func (du *DepartmentUseCase) Delete(ctx context.Context, id uint) error {
	tx := du.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	department := &entity.Department{}
	if err := du.DepartmentRepository.FindById(tx, department, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			du.Log.WithField("id", id).Error("Department not found for deletion")
			return fiber.NewError(fiber.StatusNotFound, "Department not found")
		}
		du.Log.WithError(err).Error("Error retrieving department for deletion")
		return err
	}

	if err := du.DepartmentRepository.Delete(tx, department); err != nil {
		du.Log.WithError(err).Error("Error deleting department")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		du.Log.WithError(err).Error("Error committing transaction in DeleteDepartment")
		return err
	}

	return nil
}
