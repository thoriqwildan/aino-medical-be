package http

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type DepartmentController struct {
	DepartmentUseCase *usecase.DepartmentUseCase
	Log *logrus.Logger
}

func NewDepartmentController(usecase *usecase.DepartmentUseCase, log *logrus.Logger) *DepartmentController {
	return &DepartmentController{
		DepartmentUseCase: usecase,
		Log: log,
	}
}