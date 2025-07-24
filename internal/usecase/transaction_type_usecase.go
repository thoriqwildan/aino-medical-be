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

type TransactionTypeUseCase struct {
	Repository *repository.TransactionTypeRepository
	DB *gorm.DB
	Log *logrus.Logger
	Validate *validator.Validate
}

func NewTransactionTypeUseCase(db *gorm.DB, log *logrus.Logger, tt *repository.TransactionTypeRepository, validate *validator.Validate) *TransactionTypeUseCase {
	return &TransactionTypeUseCase{
		Repository: tt,
		DB: db,
		Log: log,
		Validate: validate,
	}
}

func (ttu *TransactionTypeUseCase) Create(ctx context.Context, request *model.TransactionTypeRequest) (*model.TransactionTypeResponse, error) {
	tx := ttu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ttu.Validate.Struct(request); err != nil {
		ttu.Log.WithError(err).Error("Validation error in CreateTransactionType")
		return nil, err
	}

	if err := ttu.Repository.FindByName(tx, request.Name); err == nil {
		ttu.Log.WithField("name", request.Name).Error("Transaction type already exists")
		return nil, fiber.ErrConflict
	}

	transactionType := &entity.TransactionType{
		Name: request.Name,
	}

	if err := ttu.Repository.Create(tx, transactionType); err != nil {
		ttu.Log.WithError(err).Error("Error creating transaction type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ttu.Log.WithError(err).Error("Error committing transaction in CreateTransactionType")
		return nil, err
	}

	return converter.TransactionTypeToResponse(transactionType), nil
}

func (ttu *TransactionTypeUseCase) GetById(ctx context.Context, id int) (*model.TransactionTypeResponse, error) {
	tx := ttu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	transactionType := &entity.TransactionType{}
	if err := ttu.Repository.FindById(tx, transactionType, id); err != nil {
		ttu.Log.WithError(err).Error("Error finding transaction type by ID")
		return nil, err
	}

	return converter.TransactionTypeToResponse(transactionType), nil
}