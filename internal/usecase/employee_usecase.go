package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"gorm.io/gorm"
)

type EmployeeUseCase struct {
	Repository *repository.EmployeeRepository
	Log        *logrus.Logger
	DB         *gorm.DB
	Validate   *validator.Validate
}

func NewEmployeeUseCase(db *gorm.DB, log *logrus.Logger, er *repository.EmployeeRepository, validate *validator.Validate) *EmployeeUseCase {
	return &EmployeeUseCase{
		Repository: er,
		Log:        log,
		DB:         db,
		Validate:   validate,
	}
}