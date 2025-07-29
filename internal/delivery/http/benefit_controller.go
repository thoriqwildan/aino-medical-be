package http

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type BenefitController struct {
	UseCase *usecase.BenefitUseCase
	Log *logrus.Logger
}

func NewBenefitController(useCase *usecase.BenefitUseCase, log *logrus.Logger) *BenefitController {
	return &BenefitController{
		UseCase: useCase,
		Log: log,
	}
}