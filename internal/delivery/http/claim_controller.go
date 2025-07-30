package http

import (
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type ClaimController struct {
	UseCase *usecase.ClaimUseCase
	Log *logrus.Logger
}

func NewClaimController(useCase *usecase.ClaimUseCase, log *logrus.Logger) *ClaimController {
	return &ClaimController{
		UseCase: useCase,
		Log: log,
	}
}