package repository

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"gorm.io/gorm"
)

type TransactionTypeRepository struct {
	Repository[entity.TransactionType]
	Log *logrus.Logger
}

func NewTransactionTypeRepository(log *logrus.Logger) *TransactionTypeRepository {
	return &TransactionTypeRepository{
		Log: log,
	}
}

func (ttr *TransactionTypeRepository) FindOrCreate(db *gorm.DB, transactionType *entity.TransactionType) (*entity.TransactionType, error) {
	var result entity.TransactionType
	errTake := db.Model(entity.TransactionType{}).Where("name = ?", transactionType.Name).Take(&result).Error
	if errTake != nil && !errors.Is(errTake, gorm.ErrRecordNotFound) {
		ttr.Log.Error("Error when find transaction type in method find or create", errTake.Error())
		return nil, errTake
	}

	if errTake == nil {
		return &result, nil 
	}

	errCreate := db.Create(transactionType).Error
	if errCreate != nil {
		ttr.Log.Error("Error when create transaction type in method find or create")
		return nil, errCreate
	}

	return transactionType, nil
}

func (ttr *TransactionTypeRepository) FindByName(db *gorm.DB, name string) error {
	return db.Where("name = ?", name).First(&entity.TransactionType{}).Error
}

func (ttr *TransactionTypeRepository) Search(db *gorm.DB, request *model.PagingQuery) ([]entity.TransactionType, int64, error) {
	var transactionTypes []entity.TransactionType
	var total int64

	baseQuery := db.Model(&entity.TransactionType{})

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if request.Limit > 0 {
		baseQuery = baseQuery.Limit(request.Limit)
	}
	if request.Page > 0 {
		baseQuery = baseQuery.Offset((request.Page - 1) * request.Limit)
	}
	err := baseQuery.
		Find(&transactionTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return transactionTypes, total, nil
}
