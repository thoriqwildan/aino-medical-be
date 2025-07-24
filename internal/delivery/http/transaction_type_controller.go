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

func (c *TransactionTypeController) Get(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := c.UseCase.Get(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching transaction types")
		return err
	}

	paging := &model.PaginationPage{
		Page: query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.TransactionTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Transaction types fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

func (c *TransactionTypeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		c.Log.Error("ID is required for update")
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateTransactionTypeRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	request.ID = uint(idInt)

	response, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		if err.Error() == "record not found" {
			return fiber.NewError(fiber.StatusNotFound, "Transaction type not found")
		}
		c.Log.WithError(err).Error("Error updating transaction type")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.TransactionTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Transaction type updated successfully",
		Data: response,
	})
}

func (c *TransactionTypeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		c.Log.Error("ID is required for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := c.UseCase.Delete(ctx.Context(), idInt); err != nil {
		if err.Error() == "record not found" {
			return fiber.NewError(fiber.StatusNotFound, "Transaction type not found")
		}
		c.Log.WithError(err).Error("Error deleting transaction type")
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[any]{
		Code: fiber.StatusNoContent,
		Message: "Transaction type deleted successfully",
		Data: nil,
	})
}

