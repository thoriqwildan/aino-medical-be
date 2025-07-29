package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type BenefitUseCase struct {
	Repository *repository.BenefitRepository
	Validate *validator.Validate
	DB *gorm.DB
	Log *logrus.Logger
}

func NewBenefitUseCase(repo *repository.BenefitRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) *BenefitUseCase {
	return &BenefitUseCase{
		Repository: repo,
		Validate: validate,
		DB: db,
		Log: log,
	}
}