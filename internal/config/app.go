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

	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, userRepository, config.Validate)

	userController := http.NewUserController(userUseCase, config.Log, config.Config)

	routeConfig := route.RouteConfig{
		App: config.App,
		JWT: config.JWT,
		UserController: userController,
	}

	routeConfig.Setup()
}