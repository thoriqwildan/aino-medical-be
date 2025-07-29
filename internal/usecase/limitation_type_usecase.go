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

type LimitationTypeUseCase struct {
	Repository *repository.LimitationTypeRepository
	DB *gorm.DB
	Log *logrus.Logger
	Validate *validator.Validate
}

func NewLimitationTypeUseCase(r *repository.LimitationTypeRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) *LimitationTypeUseCase {
	return &LimitationTypeUseCase{
		Repository: r,
		DB:        db,
		Log:      log,
		Validate: validate,
	}
}

func (ltu *LimitationTypeUseCase) Create(ctx context.Context, request *model.LimitationTypeRequest) (*model.LimitationTypeResponse, error) {
	tx := ltu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ltu.Validate.Struct(request); err != nil {
		ltu.Log.WithError(err).Error("Validation error in CreateLimitationType")
		return nil, err
	}

	if err := ltu.Repository.GetByName(tx, request.Name); err == nil {
		ltu.Log.WithField("name", request.Name).Error("Limitation type already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Limitation type already exists")
	}

	limitationType := &entity.LimitationType{
		Name: request.Name,
	}

	if err := ltu.Repository.Create(tx, limitationType); err != nil {
		ltu.Log.WithError(err).Error("Error creating limitation type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ltu.Log.WithError(err).Error("Error committing transaction in CreateLimitationType")
		return nil, err
	}

	return converter.LimitationTypeToResponse(limitationType), nil
}

func (ltu *LimitationTypeUseCase) GetById(ctx context.Context, id uint) (*model.LimitationTypeResponse, error) {
	tx := ltu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	limitationType := &entity.LimitationType{}
	if err := ltu.Repository.FindById(tx, limitationType, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Limitation type not found")
	}

	if err := tx.Commit().Error; err != nil {
		ltu.Log.WithError(err).Error("Error committing transaction in GetById")
		return nil, err
	}

	return converter.LimitationTypeToResponse(limitationType), nil
}

func (ltu *LimitationTypeUseCase) GetAll(ctx context.Context, request *model.PagingQuery) ([]model.LimitationTypeResponse, int64, error) {
	tx := ltu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ltu.Validate.Struct(request); err != nil {{
		ltu.Log.WithError(err).Error("Validation error in GetAllLimitationTypes")
		return nil, 0, err
	}}

	limitationTypes, total, err := ltu.Repository.SearchLimitationTypes(tx, request)
	if err != nil {
		ltu.Log.WithError(err).Error("Error searching limitation types")
		return nil, 0, err
	}

	responses := make([]model.LimitationTypeResponse, len(limitationTypes))
	for i, lt := range limitationTypes {
		responses[i] = *converter.LimitationTypeToResponse(&lt)
	}
	return responses, total, nil
}

func (ltu *LimitationTypeUseCase) Update(ctx context.Context, request *model.UpdateLimitationTypeRequest) (*model.LimitationTypeResponse, error) {
	tx := ltu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ltu.Validate.Struct(request); err != nil {
		ltu.Log.WithError(err).Error("Validation error in UpdateLimitationType")
		return nil, err
	}

	limitationType := &entity.LimitationType{}
	if err := ltu.Repository.FindById(tx, limitationType, request.ID); err != nil {
		ltu.Log.WithError(err).Error("Error finding limitation type by ID for update")
		return nil, fiber.NewError(fiber.StatusNotFound, "Limitation type not found")
	}

	if err := ltu.Repository.GetByName(tx, request.Name); err == nil {
		ltu.Log.WithField("name", request.Name).Error("Limitation type with this name already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Limitation type with this name already exists")
	}

	limitationType.Name = request.Name

	if err := ltu.Repository.Update(tx, limitationType); err != nil {
		ltu.Log.WithError(err).Error("Error updating limitation type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ltu.Log.WithError(err).Error("Error committing transaction in UpdateLimitationType")
		return nil, err
	}

	return converter.LimitationTypeToResponse(limitationType), nil
}

func (ltu *LimitationTypeUseCase) Delete(ctx context.Context, id uint) error {
	tx := ltu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	limitationType := &entity.LimitationType{}
	if err := ltu.Repository.FindById(tx, limitationType, id); err != nil {
		ltu.Log.WithError(err).Error("Error finding limitation type by ID for deletion")
		return fiber.NewError(fiber.StatusNotFound, "Limitation type not found")
	}

	if err := ltu.Repository.Delete(tx, limitationType); err != nil {
		ltu.Log.WithError(err).Error("Error deleting limitation type")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ltu.Log.WithError(err).Error("Error committing transaction in DeleteLimitationType")
		return err
	}

	return nil
}