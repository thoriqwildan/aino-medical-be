package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type FamilyMemberUseCase struct {
	Repository *repository.FamilyMemberRepository
	DB *gorm.DB
	Validate *validator.Validate
	Log *logrus.Logger
}

func NewFamilyMemberUseCase(repo *repository.FamilyMemberRepository, db *gorm.DB, validate *validator.Validate, log *logrus.Logger) *FamilyMemberUseCase {
	return &FamilyMemberUseCase{
		Repository: repo,
		DB: db,
		Validate: validate,
		Log: log,
	}
}