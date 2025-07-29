package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/http"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/http/route"
	"github.com/thoriqwildan/aino-medical-be/internal/delivery/middleware"
	"github.com/thoriqwildan/aino-medical-be/internal/repository"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB *gorm.DB
	App *fiber.App
	Log *logrus.Logger
	Validate *validator.Validate
	Config *viper.Viper
	JWT *middleware.MiddlewareConfig
}

func Bootstrap(config *BootstrapConfig) {
	userRepository := repository.NewUserRepository(config.Log)
	transactionTypeRepository := repository.NewTransactionTypeRepository(config.Log)
	planTypeRepository := repository.NewPlanTypeRepository(config.Log)
	limitationTypeRepository := repository.NewLimitationTypeRepository(config.Log)

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, userRepository, config.Validate)
	transactionTypeUseCase := usecase.NewTransactionTypeUseCase(config.DB, config.Log, transactionTypeRepository, config.Validate)
	planTypeUseCase := usecase.NewPlanTypeUseCase(config.DB, config.Log, planTypeRepository, config.Validate)
	limitationTypeUseCase := usecase.NewLimitationTypeUseCase(limitationTypeRepository, config.DB, config.Log, config.Validate)

	userController := http.NewUserController(userUseCase, config.Log, config.Config)
	transactionTypeController := http.NewTransactionTypeController(transactionTypeUseCase, config.Log, config.Config)
	planTypeController := http.NewPlanTypeController(planTypeUseCase, config.Log, config.Config)
	limitationTypeController := http.NewLimitationTypeController(limitationTypeUseCase, config.Log, config.Config)

	routeConfig := route.RouteConfig{
		App: config.App,
		JWT: config.JWT,
		UserController: userController,
		TransactionTypeController: transactionTypeController,
		PlanTypeController: planTypeController,
		LimitationTypeController: limitationTypeController,
	}

	routeConfig.Setup()
}