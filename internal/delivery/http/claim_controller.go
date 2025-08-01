package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type ClaimController struct {
	UseCase *usecase.ClaimUseCase
	Log *logrus.Logger
}

func NewClaimController(useCase *usecase.ClaimUseCase, log *logrus.Logger) *ClaimController {
	return &ClaimController{
		UseCase: useCase,
		Log: log,
	}
}

// @Router /api/v1/claims [post]
// @Param  request body model.ClaimRequest true "Create Claim Request"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth api_key
// @Summary Create a new claim
// @Description Create a new claim with the provided details.
// @Accept json
func (c *ClaimController) CreateClaim(ctx *fiber.Ctx) error {
	request := new(model.ClaimRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create claim")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.ClaimResponse]{
		Code: 	fiber.StatusCreated,
		Message: "Claim created successfully",
		Data: 	response,
	})
}

// @Router /api/v1/claims/get-patients [get]
// @Success 200 {object} model.PatientResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth api_key
// @Summary Find patients
// @Description Find patients by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *ClaimController) GetAllPatient(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := c.UseCase.GetPatient(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching patients")
		return err
	}

	paging := &model.PaginationPage{
		Page: query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.PatientResponse]{
		Code: fiber.StatusOK,
		Message: "Patients fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

// @Router /api/v1/claims/get-benefits/{patientId} [get]
// @Param  patientId path int true "Patient ID"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth api_key
// @Summary Find benefits
// @Description Find benefits by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *ClaimController) GetAllBenefits(ctx *fiber.Ctx) error {
	id := ctx.Params("patientId")
	if id == "" {
		c.Log.Error("Patient ID is required")
		return fiber.NewError(fiber.StatusBadRequest, "Patient ID is required")
	}

	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	patientId, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid patient ID format")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid patient ID format")
	}

	responses, total, err := c.UseCase.GetBenefit(ctx.Context(), query, uint(patientId))
	if err != nil {
		c.Log.WithError(err).Error("Error fetching benefits")
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

// @Router /api/v1/claims/{id} [put]
// @Param  request body model.UpdateClaimRequest true "Update Claim Request"
// @Param id path string true "Claim ID"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth api_key
// @Summary Update a claim
// @Description Update a claim with the provided details.
// @Accept json
func (c *ClaimController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateClaimRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	request.ID = uint(idInt)

	response, err := c.UseCase.UpdateClaim(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error updating claim")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.ClaimResponse]{
		Code: fiber.StatusOK,
		Message: "Claim updated successfully",
		Data: response,
	})
}