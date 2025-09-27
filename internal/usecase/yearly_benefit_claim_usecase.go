package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/model/converter"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type YearlyBenefitClaimUsecase struct {
	Repository *repository.YearlyBenefitClaimRepository
	Validate   *validator.Validate
	DB         *gorm.DB
	Log        *logrus.Logger
}

func NewYearlyBenefitClaimUsecase(repository *repository.YearlyBenefitClaimRepository, validate *validator.Validate, DB *gorm.DB, log *logrus.Logger) *YearlyBenefitClaimUsecase {
	return &YearlyBenefitClaimUsecase{Repository: repository, Validate: validate, DB: DB, Log: log}
}

func (yu YearlyBenefitClaimUsecase) Create(ctx context.Context, request *model.YearlyBenefitClaimRequest) (*model.YearlyBenefitClaimResponse, error) {
	tx := yu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := yu.Validate.Struct(request); err != nil {
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("validate err: %v", err)
		return nil, err
	}
	yearlyClaim := entity.YearlyBenefitClaim{
		Code:        request.Code,
		YearlyClaim: request.YearlyClaim,
	}
	if err := yu.Repository.Create(tx, &yearlyClaim); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fiber.NewError(fiber.StatusConflict, fmt.Sprintf("error when create yearly benefit claim because duplicate key: %s", err.Error()))
		}
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("error when create yearly claim %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when create yearly benefit claim: %s", err.Error()))
	}

	return converter.YearlyBenefitClaimToResponse(&yearlyClaim), nil
}

func (yu YearlyBenefitClaimUsecase) GetByID(ctx context.Context, id uint) (*model.YearlyBenefitClaimResponse, error) {
	tx := yu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	var yearlyClaim entity.YearlyBenefitClaim
	if err := yu.Repository.FindById(tx, &yearlyClaim, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("yearly benefit claim not found")
			return nil, fiber.NewError(fiber.StatusNotFound, "yearly benefit claim not found")
		}
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("error when get yearly benefit claim: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when get yearly benefit claim: %s", err.Error()))
	}
	return converter.YearlyBenefitClaimToResponse(&yearlyClaim), nil
}

func (yu YearlyBenefitClaimUsecase) GetAll(ctx context.Context, request *model.YearlyBenefitClaimFilter) ([]*model.YearlyBenefitClaimResponse, int64, error) {
	tx := yu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if yearlyBenefitClaims, count, err := yu.Repository.GetAll(tx, request); err != nil {
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("error when get all yearly benefit claims: %v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when get all yearly benefit claims: %v", err))
	} else {
		return converter.YearlyBenefitClaimToResponses(yearlyBenefitClaims), count, nil
	}
}

func (yu YearlyBenefitClaimUsecase) Update(ctx context.Context, id uint, request *model.UpdateYearlyBenefitClaimRequest) (*model.YearlyBenefitClaimResponse, error) {
	tx := yu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := yu.Validate.Struct(request); err != nil {
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("validate err: %v", err)
		return nil, err
	}

	if count, err := yu.Repository.CountById(tx, id); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when counting yearly benefit claim: %v", err))
	} else if count == 0 || count < 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "no yearly benefit claim found")
	}

	yearlyClaim := entity.YearlyBenefitClaim{
		ID:          id,
		Code:        request.Code,
		YearlyClaim: request.YearlyClaim,
	}
	if err := yu.Repository.Update(tx, &yearlyClaim); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fiber.NewError(fiber.StatusConflict, fmt.Sprintf("error when update yearly benefit claim because duplicate key: %s", err.Error()))
		}
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("error when update yearly claim %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when update yearly benefit claim: %s", err.Error()))
	}

	return converter.YearlyBenefitClaimToResponse(&yearlyClaim), nil
}

func (yu YearlyBenefitClaimUsecase) Delete(ctx context.Context, id uint) error {
	tx := yu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if count, err := yu.Repository.CountById(tx, id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when counting yearly benefit claim: %v", err))
	} else if count == 0 || count < 0 {
		return fiber.NewError(fiber.StatusNotFound, "no yearly benefit claim found")
	}
	if err := yu.DB.Where("id = ?", id).Delete(entity.YearlyBenefitClaim{}).Error; err != nil {
		yu.Log.WithField("usecase", "YearlyBenefitClaimUsecase").Errorf("error when delete yearly claim %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("error when delete yearly benefit claim: %s", err.Error()))
	}

	return nil
}
