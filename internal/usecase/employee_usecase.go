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

type EmployeeUseCase struct {
	Repository *repository.EmployeeRepository
	Log        *logrus.Logger
	DB         *gorm.DB
	Validate   *validator.Validate
}

func NewEmployeeUseCase(db *gorm.DB, log *logrus.Logger, er *repository.EmployeeRepository, validate *validator.Validate) *EmployeeUseCase {
	return &EmployeeUseCase{
		Repository: er,
		Log:        log,
		DB:         db,
		Validate:   validate,
	}
}

func (eu *EmployeeUseCase) Create(ctx context.Context, request *model.EmployeeRequest) (*model.EmployeeResponse, error) {
	tx := eu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := eu.Validate.Struct(request); err != nil {
		eu.Log.WithError(err).Error("Validation error in CreateEmployee")
		return nil, err
	}

	if err := eu.Repository.GetDepartmentByID(tx, request.DepartmentID); err != nil {
		eu.Log.WithField("department_id", request.DepartmentID).Error("Department not found in CreateEmployee")
		return nil, fiber.NewError(fiber.StatusNotFound, "Department not found")
	}

	if err := eu.Repository.GetByEmail(tx, request.Email); err == nil {
		eu.Log.WithField("email", request.Email).Error("Employee with this email already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Employee with this email already exists")
	}

	employee := &entity.Employee{
		Name: request.Name,
		DepartmentID: request.DepartmentID,
		Position: request.Position,
		Email: request.Email,
		Phone: request.Phone,
		BirthDate: time.Time(request.BirthDate),
		Gender: entity.Genders(request.Gender),
		PlanTypeID: request.PlanTypeID,
		Dependence: &request.Dependences,
		BankNumber: request.BankNumber,
		JoinDate: time.Time(request.JoinDate),
		Patient: entity.Patient{
			PlanTypeID: request.PlanTypeID,
			Name: request.Name,
			BirthDate: time.Time(request.BirthDate),
			Gender: entity.Genders(request.Gender),
			FamilyMemberID: nil,
		},
	}

	if err := eu.Repository.Create(tx, employee); err != nil {
		eu.Log.WithError(err).Error("Error creating employee in CreateEmployee")
		return nil, err
	}

	if err := eu.Repository.FindById(tx, employee.ID, employee); err != nil {
		eu.Log.WithError(err).Error("Error finding employee by ID in CreateEmployee")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		eu.Log.WithError(err).Error("Error committing transaction in CreateEmployee")
		return nil, err
	}

	return converter.EmployeeToResponse(employee), nil
}

func (eu *EmployeeUseCase) GetById(ctx context.Context, id uint) (*model.EmployeeResponse, error) {
	tx := eu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	employee := &entity.Employee{}
	if err := eu.Repository.FindById(tx, id, employee); err != nil {
		if err == gorm.ErrRecordNotFound {
			eu.Log.WithField("id", id).Error("Employee not found in GetById")
			return nil, fiber.NewError(fiber.StatusNotFound, "Employee not found")
		}
		eu.Log.WithError(err).Error("Error finding employee by ID in GetById")
		return nil, err
	}

	return converter.EmployeeToResponse(employee), nil
}

func (eu *EmployeeUseCase) GetAll(ctx context.Context, request *model.PagingQuery) ([]model.EmployeeResponse, int64, error) {
	tx := eu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := eu.Validate.Struct(request); err != nil {
		eu.Log.WithError(err).Error("Validation error in GetAllEmployees")
		return nil, 0, err
	}

	employees, total, err := eu.Repository.SearchEmployees(tx, request)
	if err != nil {
		eu.Log.WithError(err).Error("Error searching employees")
		return nil, 0, err
	}

	responses := make([]model.EmployeeResponse, len(employees))
	for i, emp := range employees {
		responses[i] = *converter.EmployeeToResponse(&emp)
	}
	return responses, total, nil
}

func (eu *EmployeeUseCase) Update(ctx context.Context, request *model.UpdateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := eu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := eu.Validate.Struct(request); err != nil {
		eu.Log.WithError(err).Error("Validation error in UpdateEmployee")
		return nil, err
	}

	employee := &entity.Employee{}
	if err := eu.Repository.FindById(tx, request.ID, employee); err != nil {
		if err == gorm.ErrRecordNotFound {
			eu.Log.WithField("id", request.ID).Error("Employee not found in UpdateEmployee")
			return nil, fiber.NewError(fiber.StatusNotFound, "Employee not found")
		}
		eu.Log.WithError(err).Error("Error finding employee by ID in UpdateEmployee")
		return nil, err
	}

	if err := eu.Repository.GetDepartmentByID(tx, request.DepartmentID); err != nil {
		eu.Log.WithField("department_id", request.DepartmentID).Error("Department not found in UpdateEmployee")
		return nil, fiber.NewError(fiber.StatusNotFound, "Department not found")
	}

	eu.Log.Info("Plan Type ID:", request.PlanTypeID)

	// Update fields
	employee.Name = request.Name
	employee.Position = request.Position
	employee.Email = request.Email
	employee.Phone = request.Phone
	employee.BirthDate = time.Time(request.BirthDate)
	employee.Gender = entity.Genders(request.Gender)
	employee.PlanTypeID = request.PlanTypeID
	employee.DepartmentID = request.DepartmentID
	employee.Dependence = &request.Dependences
	employee.BankNumber = request.BankNumber
	employee.JoinDate = time.Time(request.JoinDate)

	if err := eu.Repository.Update(tx, employee); err != nil {
		eu.Log.WithError(err).Error("Error updating employee in UpdateEmployee")
		return nil, err
	}

	if err := eu.Repository.FindById(tx, employee.ID, employee); err != nil {
		eu.Log.WithError(err).Error("Error finding employee by ID after update in UpdateEmployee")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		eu.Log.WithError(err).Error("Error committing transaction in UpdateEmployee")
		return nil, err
	}

	return converter.EmployeeToResponse(employee), nil
}

func (eu *EmployeeUseCase) Delete(ctx context.Context, id uint) error {
	tx := eu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	employee := &entity.Employee{}
	if err := eu.Repository.FindById(tx, id, employee); err != nil {
		if err == gorm.ErrRecordNotFound {
			eu.Log.WithField("id", id).Error("Employee not found in Delete")
			return fiber.NewError(fiber.StatusNotFound, "Employee not found")
		}
		eu.Log.WithError(err).Error("Error finding employee by ID in Delete")
		return err
	}

	if err := eu.Repository.Delete(tx, employee); err != nil {
		eu.Log.WithError(err).Error("Error deleting employee in Delete")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		eu.Log.WithError(err).Error("Error committing transaction in Delete")
		return err
	}

	return nil
}