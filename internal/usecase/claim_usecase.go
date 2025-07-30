package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type ClaimUseCase struct {
	Repository *repository.ClaimRepository
	Log *logrus.Logger
	DB *gorm.DB
	Validate *validator.Validate
}

func NewClaimUseCase(repo *repository.ClaimRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger) *ClaimUseCase {
	return &ClaimUseCase{
		Repository: repo,
		DB: db,
		Validate: validate,
		Log: log,
	}
}