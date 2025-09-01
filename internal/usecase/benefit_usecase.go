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

type BenefitUseCase struct {
	Repository *repository.BenefitRepository
	Validate   *validator.Validate
	DB         *gorm.DB
	Log        *logrus.Logger
}

func NewBenefitUseCase(repo *repository.BenefitRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) *BenefitUseCase {
	return &BenefitUseCase{
		Repository: repo,
		Validate:   validate,
		DB:         db,
		Log:        log,
	}
}

func (bu *BenefitUseCase) Create(ctx context.Context, request *model.CreateBenefitRequest) (*model.BenefitResponse, error) {
	tx := bu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := bu.Validate.Struct(request); err != nil {
		bu.Log.WithError(err).Error("Validation error in CreateBenefit")
		return nil, err
	}

	if err := bu.Repository.GetByCode(tx, request.Code); err == nil {
		bu.Log.WithField("code", request.Code).Error("Benefit with this code already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Benefit with this code already exists")
	}

	benefit := &entity.Benefit{
		Name:             request.Name,
		PlanTypeID:       request.PlanTypeID,
		Detail:           request.Detail,
		Code:             request.Code,
		LimitationTypeID: request.LimitationTypeID,
		Plafond:          &request.Plafond,
		YearlyMax:        &request.YearlyMax,
	}
	if err := bu.Repository.Create(tx, benefit); err != nil {
		bu.Log.WithError(err).Error("Error creating benefit")
		return nil, err
	}
	if err := bu.Repository.GetById(tx, benefit.ID, benefit); err != nil {
		bu.Log.WithError(err).Error("Error retrieving created benefit")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Error retrieving created benefit")
	}

	if err := tx.Commit().Error; err != nil {
		bu.Log.WithError(err).Error("Error committing transaction in CreateBenefit")
		return nil, err
	}

	return converter.BenefitToResponse(benefit), nil
}

func (bu *BenefitUseCase) GetById(ctx context.Context, id uint) (*model.BenefitResponse, error) {
	tx := bu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	benefit := &entity.Benefit{}
	if err := bu.Repository.GetById(tx, id, benefit); err != nil {
		bu.Log.WithError(err).Error("Error finding benefit by ID")
		return nil, fiber.NewError(fiber.StatusNotFound, "Benefit not found")
	}

	if err := tx.Commit().Error; err != nil {
		bu.Log.WithError(err).Error("Error committing transaction in GetById")
		return nil, err
	}

	return converter.BenefitToResponse(benefit), nil
}

func (bu *BenefitUseCase) GetAll(ctx context.Context, request *model.SearchPagingQuery) ([]model.BenefitResponse, int64, error) {
	tx := bu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := bu.Validate.Struct(request); err != nil {
		bu.Log.WithError(err).Error("Validation error in GetAllBenefits")
		return nil, 0, err
	}

	benefits, total, err := bu.Repository.SearchBenefits(tx, request)
	if err != nil {
		bu.Log.WithError(err).Error("Error searching benefits")
		return nil, 0, err
	}

	responses := make([]model.BenefitResponse, len(benefits))
	for i, b := range benefits {
		responses[i] = *converter.BenefitToResponse(&b)
	}
	return responses, total, nil
}

func (bu *BenefitUseCase) Update(ctx context.Context, request *model.UpdateBenefitRequest) (*model.BenefitResponse, error) {
	tx := bu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := bu.Validate.Struct(request); err != nil {
		bu.Log.WithError(err).Error("Validation error in UpdateBenefit")
		return nil, err
	}

	benefit := &entity.Benefit{
		ID:               request.ID,
		Name:             request.Name,
		PlanTypeID:       request.PlanTypeID,
		Detail:           request.Detail,
		Code:             request.Code,
		LimitationTypeID: request.LimitationTypeID,
		Plafond:          &request.Plafond,
		YearlyMax:        &request.YearlyMax,
	}

	if err := bu.Repository.Update(tx, benefit); err != nil {
		bu.Log.WithError(err).Error("Error updating benefit")
		return nil, err
	}

	if err := bu.Repository.GetById(tx, benefit.ID, benefit); err != nil {
		bu.Log.WithError(err).Error("Error retrieving updated benefit")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		bu.Log.WithError(err).Error("Error committing transaction in UpdateBenefit")
		return nil, err
	}

	return converter.BenefitToResponse(benefit), nil
}

func (bu *BenefitUseCase) Delete(ctx context.Context, id uint) error {
	tx := bu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	benefit := &entity.Benefit{}
	if err := bu.Repository.GetById(tx, id, benefit); err != nil {
		bu.Log.WithError(err).Error("Error finding benefit by ID for deletion")
		return fiber.NewError(fiber.StatusNotFound, "Benefit not found")
	}

	if err := bu.Repository.Delete(tx, benefit); err != nil {
		bu.Log.WithError(err).Error("Error deleting benefit")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		bu.Log.WithError(err).Error("Error committing transaction in DeleteBenefit")
		return err
	}

	return nil
}
