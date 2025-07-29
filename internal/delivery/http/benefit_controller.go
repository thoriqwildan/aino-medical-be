package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
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

func (c *BenefitController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateBenefitRequest)

	ctx.BodyParser(&request)

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating benefit")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.BenefitResponse]{
		Code: fiber.StatusCreated,
		Message: "Benefit created successfully",
		Data: response,
	})
}

func (c *BenefitController) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.GetById(ctx.Context(), idUint)
	if err != nil {
		c.Log.WithError(err).Error("Error retrieving benefit by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.BenefitResponse]{
		Code: fiber.StatusOK,
		Message: "Benefit retrieved successfully",
		Data: response,
	})
}

func (c *BenefitController) GetAll(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := c.UseCase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching plan types")
		return err
	}

	paging := &model.PaginationPage{
		Page: query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.BenefitResponse]{
		Code: fiber.StatusOK,
		Message: "Benefits fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

func (c *BenefitController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateBenefitRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	request.ID = uint(idInt)

	response, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error updating limitation type")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.BenefitResponse]{
		Code: fiber.StatusOK,
		Message: "Benefit updated successfully",
		Data: response,
	})
}

func (c *BenefitController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	err = c.UseCase.Delete(ctx.Context(), idUint)
	if err != nil {
		c.Log.WithError(err).Error("Error deleting benefit")
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[any]{
		Code: fiber.StatusNoContent,
		Message: "Benefit deleted successfully",
		Data: nil,
	})
}