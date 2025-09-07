package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type YearlyBenefitClaimController struct {
	UseCase *usecase.YearlyBenefitClaimUsecase
	Log     *logrus.Logger
}

func NewYearlyBenefitClaimController(useCase *usecase.YearlyBenefitClaimUsecase, log *logrus.Logger) *YearlyBenefitClaimController {
	return &YearlyBenefitClaimController{UseCase: useCase, Log: log}
}

// @Router /api/v1/yearly-claims [post]
// @Param  request body model.YearlyBenefitClaimRequest true "Create Plan Type Request"
// @Success 200 {object} model.YearlyBenefitClaimWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Yearly Benefit Claim
// @Security    BearerAuth bearer
// @Summary Create a new yearly benefit claim
// @Description Create a new yearly benefit claim with the provided details.
// @Accept json
func (c *YearlyBenefitClaimController) Create(ctx *fiber.Ctx) error {
	request := new(model.YearlyBenefitClaimRequest)

	err := ctx.BodyParser(&request)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error creating yearly benefit claim")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.YearlyBenefitClaimResponse]{
		Code:    fiber.StatusCreated,
		Message: "Yearly benefit claim created successfully",
		Data:    response,
	})
}

// @Router /api/v1/yearly-claims/{id} [get]
// @Param  id path int true "Yearly Benefit Claim ID"
// @Success 200 {object} model.YearlyBenefitClaimWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Yearly Benefit Claim
// @Security    BearerAuth bearer
// @Summary Get a yearly benefit claim ID
// @Description Get a yearly benefit claim its ID.
// @Accept json
func (c *YearlyBenefitClaimController) GetById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	if idStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	// Convert idStr to uint

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.GetByID(ctx.Context(), uint(id))
	if err != nil {
		c.Log.WithError(err).Error("Error getting yearly benefit claim by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.YearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Yearly benefit claim retrieved successfully",
		Data:    response,
	})
}

// @Router /api/v1/yearly-claims [get]
// @Success 200 {object} model.YearlyBenefitClaimWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Yearly Benefit Claim
// @Security    BearerAuth bearer
// @Summary Find yearly benefit claim
// @Description Find yearly benefit claim by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Param   code query    string               false       "Code for find the yearly benefit claim controller" default
// @Accept json
func (c *YearlyBenefitClaimController) Get(ctx *fiber.Ctx) error {
	query := &model.YearlyBenefitClaimFilter{
		Page:  ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
		Code:  ctx.Query("code"),
	}

	responses, total, err := c.UseCase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching yearly benefit claims")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.YearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Yearly benefit claim fetched successfully",
		Data:    &responses,
		Meta:    paging,
	})
}

// @Router /api/v1/yearly-claims/{id} [put]
// @Param  request body model.UpdateYearlyBenefitClaimRequest true "Update Yearly Benefit Claim"
// @Param id path string true "Yearly Benefit Claim ID"
// @Success 200 {object} model.YearlyBenefitClaimWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Yearly Benefit Claim
// @Security    BearerAuth bearer
// @Summary Update a yearly benefit claim
// @Description Update a yearly benefit claim with the provided details.
// @Accept json
func (c *YearlyBenefitClaimController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateYearlyBenefitClaimRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := c.UseCase.Update(ctx.Context(), uint(idInt), request)
	if err != nil {
		c.Log.WithError(err).Error("Error updating yearly benefit claim")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.YearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Yearly benefit claim updated successfully",
		Data:    response,
	})
}

// @Router /api/v1/yearly-claims/{id} [delete]
// @Param id path string true "Yearly Benefit Claim ID"
// @Success 200 {object} model.YearlyBenefitClaimWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Yearly Benefit Claim
// @Security    BearerAuth bearer
// @Summary Delete a yearly benefit claim
// @Description Delete a yearly benefit claim with the provided details.
// @Accept json
func (c *YearlyBenefitClaimController) Delete(ctx *fiber.Ctx) error {
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
		c.Log.WithError(err).Error("Error deleting yearly benefit claim by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "Yearly benefit claim deleted successfully",
	})
}
