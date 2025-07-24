package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
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