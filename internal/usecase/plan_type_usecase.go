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

type PlanTypeUseCase struct {
	Repository *repository.PlanTypeRepository 
	DB *gorm.DB
	Log *logrus.Logger
	Validate *validator.Validate
}

func NewPlanTypeUseCase(db *gorm.DB, log *logrus.Logger, pr *repository.PlanTypeRepository, validate *validator.Validate) *PlanTypeUseCase {
	return &PlanTypeUseCase{
		Repository: pr,
		DB:        db,
		Log:      log,
		Validate: validate,
	}
}

func (ptu *PlanTypeUseCase) Create(ctx context.Context, request *model.PlanTypeRequest) (*model.PlanTypeResponse, error) {
	tx := ptu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ptu.Validate.Struct(request); err != nil {
		ptu.Log.WithError(err).Error("Validation error in CreatePlanType")
		return nil, err
	}

	if err := ptu.Repository.FindByName(tx, request.Name); err == nil {
		ptu.Log.WithField("name", request.Name).Error("Plan type already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Plan type already exists")
	}

	planType := &entity.PlanType{
		Name:        request.Name,
		Description: &request.Description,
	}

	if err := ptu.Repository.Create(tx, planType); err != nil {
		ptu.Log.WithError(err).Error("Error creating plan type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ptu.Log.WithError(err).Error("Error committing transaction in CreatePlanType")
		return nil, err
	}

	return converter.PlanTypeToResponse(planType), nil
}

func (ptu *PlanTypeUseCase) GetById(ctx context.Context, id uint) (*model.PlanTypeResponse, error) {
	tx := ptu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	planType := &entity.PlanType{}
	if err := ptu.Repository.FindById(tx, planType, id); err != nil {
		ptu.Log.WithError(err).Error("Error finding plan type by ID")
		return nil, fiber.NewError(fiber.StatusNotFound, "Plan type not found")
	}

	if err := tx.Commit().Error; err != nil {
		ptu.Log.WithError(err).Error("Error committing transaction in GetById")
		return nil, err
	}

	return converter.PlanTypeToResponse(planType), nil
}

func (ptu *PlanTypeUseCase) Get(ctx context.Context, request *model.PagingQuery) ([]model.PlanTypeResponse, int64, error) {
	tx := ptu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ptu.Validate.Struct(request); err != nil {
		ptu.Log.WithError(err).Error("Validation error in GetPlanTypes")
		return nil, 0, err
	}

	planTypes, total, err := ptu.Repository.SearchPlanTypes(tx, request)
	if err != nil {
		ptu.Log.WithError(err).Error("Error searching plan types")
		return nil, 0, err
	}

	responses := make([]model.PlanTypeResponse, len(planTypes))
	for i, planType := range planTypes {
		responses[i] = *converter.PlanTypeToResponse(&planType)
	}
	return responses, total, nil
}

func (ptu *PlanTypeUseCase) Update(ctx context.Context, request *model.UpdatePlanTypeRequest) (*model.PlanTypeResponse, error) {
	tx := ptu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ptu.Validate.Struct(request); err != nil {
		ptu.Log.WithError(err).Error("Validation error in UpdatePlanType")
		return nil, err
	}

	planType := &entity.PlanType{}
	if err := ptu.Repository.FindById(tx, planType, request.ID); err != nil {
		ptu.Log.WithError(err).Error("Error finding plan type by ID for update")
		return nil, fiber.NewError(fiber.StatusNotFound, "Plan type not found")
	}

	if err := ptu.Repository.FindByName(tx, request.Name); err == nil {
		ptu.Log.WithField("name", request.Name).Error("Plan type with this name already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Plan type with this name already exists")
	}

	planType.Name = request.Name
	planType.Description = &request.Description

	if err := ptu.Repository.Update(tx, planType); err != nil {
		ptu.Log.WithError(err).Error("Error updating plan type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ptu.Log.WithError(err).Error("Error committing transaction in UpdatePlanType")
		return nil, err
	}

	return converter.PlanTypeToResponse(planType), nil
} 

func (ptu *PlanTypeUseCase) Delete(ctx context.Context, id uint) error {
	tx := ptu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	planType := &entity.PlanType{}
	if err := ptu.Repository.FindById(tx, planType, id); err != nil {
		ptu.Log.WithError(err).Error("Error finding plan type by ID for deletion")
		return fiber.NewError(fiber.StatusNotFound, "Plan type not found")
	}

	if err := ptu.Repository.Delete(tx, planType); err != nil {
		ptu.Log.WithError(err).Error("Error deleting plan type")
		return fiber.NewError(fiber.StatusNotFound, "Plan type not found")
	}

	if err := tx.Commit().Error; err != nil {
		ptu.Log.WithError(err).Error("Error committing transaction in DeletePlanType")
		return err
	}

	return nil
}