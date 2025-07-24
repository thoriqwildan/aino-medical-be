package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type TransactionTypeController struct {
	UseCase *usecase.TransactionTypeUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewTransactionTypeController(useCase *usecase.TransactionTypeUseCase, log *logrus.Logger, config *viper.Viper) *TransactionTypeController {
	return &TransactionTypeController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}

func (c *TransactionTypeController) Create(ctx *fiber.Ctx) error {
	request := new(model.TransactionTypeRequest)

	ctx.BodyParser(request)

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating transaction type")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.TransactionTypeResponse]{
		Code: fiber.StatusCreated,
		Message: "Transaction type created successfully",
		Data: response,
	})
}

func (c *TransactionTypeController) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	// Convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.GetById(ctx.Context(), idInt)
	if err != nil {
		if err.Error() == "record not found" {
			return fiber.NewError(fiber.StatusNotFound, "Transaction type not found")
		}

		c.Log.WithError(err).Error("Error fetching transaction type by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.TransactionTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Transaction type fetched successfully",
		Data: response,
	})
}
