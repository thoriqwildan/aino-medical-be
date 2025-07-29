package http

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type EmployeeController struct {
	UseCase *usecase.EmployeeUseCase
	Log *logrus.Logger
}

func NewEmployeeController(useCase *usecase.EmployeeUseCase, log *logrus.Logger) *EmployeeController {
	return &EmployeeController{
		UseCase: useCase,
		Log: log,
	}
}