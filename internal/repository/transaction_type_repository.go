package repository

import (
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

	err := baseQuery.
		Offset((request.Page - 1) * request.Limit).
		Limit(request.Limit).
		Find(&transactionTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return transactionTypes, total, nil
}
