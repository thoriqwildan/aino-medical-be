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

// @Router /api/v1/plan-types [post]
// @Param  request body model.PlanTypeRequest true "Create Plan Type Request"
// @Success 200 {object} model.PlanTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Create a new plan type
// @Description Create a new plan type with the provided details.
// @Accept json
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

// @Router /api/v1/plan-types/{id} [get]
// @Param  id path int true "Plan Type ID"
// @Success 200 {object} model.PlanTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Get a plan type by ID
// @Description Get a plan type by its ID.
// @Accept json
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

// @Router /api/v1/plan-types [get]
// @Success 200 {object} model.PlanTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Find plan types
// @Description Find plan types by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
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

// @Router /api/v1/plan-types/{id} [put]
// @Param  request body model.UpdatePlanTypeRequest true "Update Plan Type Request"
// @Param id path string true "Plan Type ID"
// @Success 200 {object} model.PlanTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Update a plan type
// @Description Update a plan type with the provided details.
// @Accept json
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

// @Router /api/v1/plan-types/{id} [delete]
// @Param id path string true "Plan Type ID"
// @Success 200 {object} model.PlanTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Delete a plan type
// @Description Delete a plan type with the provided details.
// @Accept json
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