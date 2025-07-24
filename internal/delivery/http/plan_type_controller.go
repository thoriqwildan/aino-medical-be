package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type PlanTypeController struct {
	UseCase *usecase.PlanTypeUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewPlanTypeController(useCase *usecase.PlanTypeUseCase, log *logrus.Logger, config *viper.Viper) *PlanTypeController {
	return &PlanTypeController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}

func (c *PlanTypeController) Create(ctx *fiber.Ctx) error {
	request := new(model.PlanTypeRequest)

	ctx.BodyParser(&request)

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating plan type")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.PlanTypeResponse]{
		Code: fiber.StatusCreated,
		Message: "Plan type created successfully",
		Data: response,
	})
}

func (c *PlanTypeController) GetById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	// Convert idStr to uint
	var id uint
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.GetById(ctx.Context(), id)
	if err != nil {
		c.Log.WithError(err).Error("Error getting plan type by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PlanTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Plan type retrieved successfully",
		Data: response,
	})
}

func (c *PlanTypeController) Get(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := c.UseCase.Get(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching plan types")
		return err
	}

	paging := &model.PaginationPage{
		Page: query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.PlanTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Plan types fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

func (c *PlanTypeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdatePlanTypeRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	request.ID = uint(idInt)

	response, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error updating plan type")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PlanTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Plan type updated successfully",
		Data: response,
	})
}

func (c *PlanTypeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}
	
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := c.UseCase.Delete(ctx.Context(), uint(idInt)); err != nil {
		c.Log.WithError(err).Error("Error deleting plan type")
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[any]{
		Code: fiber.StatusNoContent,
		Message: "Plan type deleted successfully",
	})
}