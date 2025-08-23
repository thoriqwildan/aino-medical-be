package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/entity"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type ClaimController struct {
	UseCase *usecase.ClaimUseCase
	Log     *logrus.Logger
}

func NewClaimController(useCase *usecase.ClaimUseCase, log *logrus.Logger) *ClaimController {
	return &ClaimController{
		UseCase: useCase,
		Log:     log,
	}
}

// @Router /api/v1/claims [post]
// @Param  request body model.ClaimRequest true "Create Claim Request"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
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
		Code:    fiber.StatusCreated,
		Message: "Claim created successfully",
		Data:    response,
	})
}

// @Router /api/v1/claims/get-patients [get]
// @Success 200 {object} model.PatientResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
// @Summary Find patients
// @Description Find patients by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *ClaimController) GetAllPatient(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page:  ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := c.UseCase.GetPatient(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching patients")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.PatientResponse]{
		Code:    fiber.StatusOK,
		Message: "Patients fetched successfully",
		Data:    &responses,
		Meta:    paging,
	})
}

// @Router /api/v1/claims/get-benefits/{patientId} [get]
// @Param  patientId path int true "Patient ID"
// @Success 200 {object} model.BenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
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
		Page:  ctx.QueryInt("page", 1),
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

// @Router /api/v1/claims/{id} [put]
// @Param  request body model.UpdateClaimRequest true "Update Claim Request"
// @Param id path string true "Claim ID"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
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
		Code:    fiber.StatusOK,
		Message: "Claim updated successfully",
		Data:    response,
	})
}

// @Router /api/v1/claims/{id} [get]
// @Param  id path int true "Claim ID"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
// @Summary Get a claim by ID
// @Description Get a claim by its ID.
// @Accept json
func (c *ClaimController) GetById(ctx *fiber.Ctx) error {
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

	response, err := c.UseCase.GetClaim(ctx.Context(), idUint)
	if err != nil {
		c.Log.WithError(err).Error("Error getting claim by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.ClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Department retrieved successfully",
		Data:    response,
	})
}

// @Router /api/v1/claims/{id} [delete]
// @Param id path string true "Claim ID"
// @Success 200 {object} model.ClaimResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
// @Summary Delete a claim
// @Description Delete a claim with the provided details.
// @Accept json
func (c *ClaimController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	err = c.UseCase.DeleteClaim(ctx.Context(), uint(idInt))
	if err != nil {
		c.Log.WithError(err).Error("Error deleting claim")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "Claim deleted successfully",
	})
}

// @Router /api/v1/claims [get]
// @Success 200 {object} model.ClaimResponseListWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Claims
// @Security    BearerAuth bearer
// @Summary Find claims
// @Description Find claims by their attributes.
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param date_from query string false "Start date for filtering in YYYY-MM-DD format"
// @Param date_to query string false "End date for filtering in YYYY-MM-DD format"
// @Param department query string false "Department name for filtering"
// @Param transaction_type query string false "Transaction type name for filtering"
// @Param sla_status query string false "SLA status for filtering (e.g., meet, overdue)"
// @Param claim_status query string false "Claim status for filtering (e.g., On Plafond, Over Plafond)"
// @Param transaction_status query string false "Transaction status for filtering (e.g., Successful, Pending, Failed)"
// @Accept json
func (c *ClaimController) GetAll(ctx *fiber.Ctx) error {
	transactionStatusStr := ctx.Query("transaction_status")
	var transactionStatus entity.TransactionStatus
	if transactionStatusStr != "" {
		transactionStatus = entity.TransactionStatus(transactionStatusStr)
	}

	query := &model.ClaimFilterQuery{
		Page:              ctx.QueryInt("page", 1),
		Limit:             ctx.QueryInt("limit", 10),
		DateFrom:          ctx.Query("date_from"),
		DateTo:            ctx.Query("date_to"),
		TransactionStatus: transactionStatus,
		Department:        ctx.Query("department"),
		TransactionType:   ctx.Query("transaction_type"),
		SLAStatus:         entity.SLA(ctx.Query("sla_status")),
		ClaimStatus:       entity.ClaimStatus(ctx.Query("claim_status")),
	}

	responses, total, err := c.UseCase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("Error fetching claims")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.ClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Claims fetched successfully",
		Data:    &responses,
		Meta:    paging,
	})
}
