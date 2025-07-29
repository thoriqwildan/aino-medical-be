package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type DepartmentUseCase struct {
	DepartmentRepository *repository.DepartmentRepository
	Validate *validator.Validate
	Log *logrus.Logger
	DB *gorm.DB
}

func NewDepartmentUseCase(repo *repository.DepartmentRepository, db *gorm.DB, log *logrus.Logger, validate *validator.Validate) *DepartmentUseCase {
	return &DepartmentUseCase{
		DepartmentRepository: repo,
		DB:                   db,
		Log:                  log,
		Validate:             validate,
	}
}