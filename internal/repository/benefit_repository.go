package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
)

type BenefitRepository struct {
	Repository[entity.Benefit]
	Log *logrus.Logger
}

func NewBenefitRepository(log *logrus.Logger) *BenefitRepository {
	return &BenefitRepository{
		Log: log,
	}
}