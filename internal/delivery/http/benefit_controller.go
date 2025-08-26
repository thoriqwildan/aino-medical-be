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
	Log     *logrus.Logger
}

func NewBenefitController(useCase *usecase.BenefitUseCase, log *logrus.Logger) *BenefitController {
	return &BenefitController{
		UseCase: useCase,
		Log:     log,
	}
}

// @Router /api/v1/benefits [post]
// @Param  request body model.CreateBenefitRequest true "Create Benefit Request"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Benefit Types
// @Security    BearerAuth bearer
// @Summary Create a new benefit type
// @Description Create a new benefit type with the provided details.
// @Accept json
func (c *BenefitController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateBenefitRequest)

	ctx.BodyParser(&request)

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating benefit")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.BenefitResponse]{
		Code:    fiber.StatusCreated,
		Message: "Benefit created successfully",
		Data:    response,
	})
}

// @Router /api/v1/benefits/{id} [get]
// @Param  id path int true "Benefit ID"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Benefit Types
// @Security    BearerAuth bearer
// @Summary Get a benefit type by ID
// @Description Get a benefit type by its ID.
// @Accept json
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
		Code:    fiber.StatusOK,
		Message: "Benefit retrieved successfully",
		Data:    response,
	})
}

// @Router /api/v1/benefits [get]
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Benefit Types
// @Security    BearerAuth bearer
// @Summary Find benefit types
// @Description Find benefit types by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *BenefitController) GetAll(ctx *fiber.Ctx) error {
	query := &model.SearchPagingQuery{
		SearchValue: ctx.Query("search_value"),
		Page:        ctx.QueryInt("page"),
		Limit:       ctx.QueryInt("limit"),
	}

	responses, total, err := c.UseCase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching plan types")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.BenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Benefits fetched successfully",
		Data:    &responses,
		Meta:    paging,
	})
}

// @Router /api/v1/benefits/{id} [put]
// @Param  request body model.UpdateBenefitRequest true "Update Benefit Request"
// @Param id path string true "Benefit ID"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Benefit Types
// @Security    BearerAuth bearer
// @Summary Update a benefit type
// @Description Update a benefit type with the provided details.
// @Accept json
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
		Code:    fiber.StatusOK,
		Message: "Benefit updated successfully",
		Data:    response,
	})
}

// @Router /api/v1/benefits/{id} [delete]
// @Param id path string true "Benefit ID"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Benefit Types
// @Security    BearerAuth bearer
// @Summary Delete a benefit type
// @Description Delete a benefit type with the provided details.
// @Accept json
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "Benefit deleted successfully",
		Data:    nil,
	})
}
