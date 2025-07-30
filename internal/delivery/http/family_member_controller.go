package http

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type FamilyMemberController struct {
	UseCase *usecase.FamilyMemberUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewFamilyMemberController(useCase *usecase.FamilyMemberUseCase, log *logrus.Logger, config *viper.Viper) *FamilyMemberController {
	return &FamilyMemberController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}