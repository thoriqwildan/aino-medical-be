package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/helper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
)

func NewFiber(viper *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Aino Medical API",
		Prefork: viper.GetBool("WEB_PREFORK"),
		ErrorHandler: NewErrorHandler(),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		if e, ok := err.(*fiber.Error); ok {
			return ctx.Status(e.Code).JSON(&model.WebResponse[any]{
				Code: e.Code,
				Message: e.Message,
				Errors: err.Error(),
			})
		} else if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorsMap := helper.TranslateErrorMessage(validationErrors)
			return ctx.Status(fiber.StatusBadRequest).JSON(model.WebResponse[any]{
				Code: fiber.StatusBadRequest,
				Message: "Validation Error",
				Errors: errorsMap,
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(model.WebResponse[any]{
			Code: fiber.StatusInternalServerError,
				Message: "Internal Server Error",
		})
	}
}

