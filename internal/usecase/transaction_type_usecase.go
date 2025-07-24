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

func (ttu *TransactionTypeUseCase) Get(ctx context.Context, request *model.PagingQuery) ([]model.TransactionTypeResponse, int64, error) {
	tx := ttu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ttu.Validate.Struct(request); err != nil {
		ttu.Log.WithError(err).Error("Validation error in GetTransactionTypes")
		return nil, 0, err
	}

	transactionTypes, total, err := ttu.Repository.Search(tx, request)
	if err != nil {
		ttu.Log.WithError(err).Error("Error searching transaction types")
		return nil, 0, err
	}

	responses := make([]model.TransactionTypeResponse, len(transactionTypes))
	for i, transactionType := range transactionTypes {
		responses[i] = *converter.TransactionTypeToResponse(&transactionType)
	}
	return responses, total, nil
}

func (ttu *TransactionTypeUseCase) Update(ctx context.Context, request *model.UpdateTransactionTypeRequest) (*model.TransactionTypeResponse, error) {
	tx := ttu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := ttu.Validate.Struct(request); err != nil {
		ttu.Log.WithError(err).Error("Validation error in UpdateTransactionType")
		return nil, err
	}

	if err := ttu.Repository.FindByName(tx, request.Name); err == nil {
		ttu.Log.WithField("name", request.Name).Error("Transaction already exists")
		return nil, fiber.NewError(fiber.StatusConflict, "Transaction type already exists")
	}

	transactionType := &entity.TransactionType{
		ID:   request.ID,
		Name: request.Name,
	}

	if err := ttu.Repository.Update(tx, transactionType); err != nil {
		ttu.Log.WithError(err).Error("Error updating transaction type")
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ttu.Log.WithError(err).Error("Error committing transaction in UpdateTransactionType")
		return nil, err
	}

	return converter.TransactionTypeToResponse(transactionType), nil
}

func (ttu *TransactionTypeUseCase) Delete(ctx context.Context, id int) error {
	tx := ttu.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	transactionType := &entity.TransactionType{}
	if err := ttu.Repository.FindById(tx, transactionType, id); err != nil {
		ttu.Log.WithError(err).Error("Error finding transaction type by ID for deletion")
		return err
	}

	if err := ttu.Repository.Delete(tx, transactionType); err != nil {
		ttu.Log.WithError(err).Error("Error deleting transaction type")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		ttu.Log.WithError(err).Error("Error committing transaction in DeleteTransactionType")
		return err
	}

	return nil
}