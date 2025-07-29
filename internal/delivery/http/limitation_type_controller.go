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

type LimitationTypeController struct {
	UseCase *usecase.LimitationTypeUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewLimitationTypeController(useCase *usecase.LimitationTypeUseCase, log *logrus.Logger, config *viper.Viper) *LimitationTypeController {
	return &LimitationTypeController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}

// @Router /api/v1/limitation-types [post]
// @Param  request body model.LimitationTypeRequest true "Create Limitation Type Request"
// @Success 200 {object} model.LimitationTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Plan Types
// @Security    BearerAuth api_key
// @Summary Create a new limitation type
// @Description Create a new limitation type with the provided details.
// @Accept json
func (c *LimitationTypeController) Create(ctx *fiber.Ctx) error {
	request := new(model.LimitationTypeRequest)

	ctx.BodyParser(&request)

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating limitation type")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.LimitationTypeResponse]{
		Code: fiber.StatusCreated,
		Message: "Limitation type created successfully",
		Data: response,
	})
}

// @Router /api/v1/limitation-types/{id} [get]
// @Param  id path int true "Limitation Type ID"
// @Success 200 {object} model.LimitationTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Limitation Types
// @Security    BearerAuth api_key
// @Summary Get a limitation type by ID
// @Description Get a limitation type by its ID.
// @Accept json
func (c *LimitationTypeController) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	// Convert id to uint
	var idUint uint
	_, err := fmt.Sscanf(id, "%d", &idUint)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.GetById(ctx.Context(), idUint)
	if err != nil {
		c.Log.WithError(err).Error("Error getting limitation type by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.LimitationTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Limitation type retrieved successfully",
		Data: response,
	})
}

// @Router /api/v1/limitation-types [get]
// @Success 200 {object} model.LimitationTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Limitation Types
// @Security    BearerAuth api_key
// @Summary Find limitation types
// @Description Find limitation types by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *LimitationTypeController) GetAll(ctx *fiber.Ctx) error {
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.LimitationTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Limitation types fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

// @Router /api/v1/limitation-types/{id} [put]
// @Param  request body model.UpdateLimitationTypeRequest true "Update Limitation Type Request"
// @Param id path string true "Limitation Type ID"
// @Success 200 {object} model.LimitationTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Limitation Types
// @Security    BearerAuth api_key
// @Summary Update a limitation type
// @Description Update a limitation type with the provided details.
// @Accept json
func (c *LimitationTypeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateLimitationTypeRequest)
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.LimitationTypeResponse]{
		Code: fiber.StatusOK,
		Message: "Limitation type updated successfully",
		Data: response,
	})
}

// @Router /api/v1/limitation-types/{id} [delete]
// @Param id path string true "Limitation Type ID"
// @Success 200 {object} model.LimitationTypeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Limitation Types
// @Security    BearerAuth api_key
// @Summary Delete a limitation type
// @Description Delete a limitation type with the provided details.
// @Accept json
func (c *LimitationTypeController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	if err := c.UseCase.Delete(ctx.Context(), uint(idUint)); err != nil {
		c.Log.WithError(err).Error("Error deleting limitation type")
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[any]{
		Code: fiber.StatusNoContent,
		Message: "Limitation type deleted successfully",
	})
}