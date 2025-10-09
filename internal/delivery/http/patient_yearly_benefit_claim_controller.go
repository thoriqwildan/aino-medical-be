package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type PatientYearlyBenefitClaimController struct {
	Log     *logrus.Logger
	Usecase *usecase.PatientYearlyBenefitClaimUsecase
}

func NewPatientYearlyBenefitClaimController(uc *usecase.PatientYearlyBenefitClaimUsecase, log *logrus.Logger) *PatientYearlyBenefitClaimController {
	return &PatientYearlyBenefitClaimController{
		Log:     log,
		Usecase: uc,
	}
}

func (c *PatientYearlyBenefitClaimController) parseIDs(patientIDStr, yearlyBenefitClaimIDStr string) (uint, uint, error) {
	pid, err := strconv.Atoi(patientIDStr)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusBadRequest, "invalid patientId format")
	}
	yid, err := strconv.Atoi(yearlyBenefitClaimIDStr)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusBadRequest, "invalid yearlyBenefitClaimId format")
	}
	return uint(pid), uint(yid), nil
}

// ========================= Handlers =========================

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/yearly-claims/{yearlyBenefitClaimId} [post]
// @Param  patientId path int true "Patient ID"
// @Param  yearlyBenefitClaimId path int true "Yearly Benefit Claim ID"
// @Success 201 {object} model.PatientYearlyBenefitClaimResponse
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Create or find patient yearly benefit claim
// @Accept json
func (c *PatientYearlyBenefitClaimController) Create(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	ybcIDStr := ctx.Params("yearlyBenefitClaimId")
	if patientIDStr == "" || ybcIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and yearlyBenefitClaimId are required")
	}
	pid, yid, err := c.parseIDs(patientIDStr, ybcIDStr)
	if err != nil {
		return err
	}

	req := &model.PatientYearlyBenefitClaimParams{
		PatientID:            pid,
		YearlyBenefitClaimID: yid,
	}

	res, err := c.Usecase.Create(ctx.Context(), req)
	if err != nil {
		c.Log.WithError(err).Error("Create patient yearly benefit claim failed")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.PatientYearlyBenefitClaimResponse]{
		Code:    fiber.StatusCreated,
		Message: "Patient yearly benefit claim created successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/yearly-claims/{yearlyBenefitClaimId} [get]
// @Param  patientId path int true "Patient ID"
// @Param  yearlyBenefitClaimId path int true "Yearly Benefit Claim ID"
// @Success 200 {object} model.PatientYearlyBenefitClaimResponse
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Get patient yearly benefit claim by patientId & yearlyBenefitClaimId
// @Accept json
func (c *PatientYearlyBenefitClaimController) GetByPatientYearlyBenefitClaimID(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	ybcIDStr := ctx.Params("yearlyBenefitClaimId")
	if patientIDStr == "" || ybcIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and yearlyBenefitClaimId are required")
	}
	pid, yid, err := c.parseIDs(patientIDStr, ybcIDStr)
	if err != nil {
		return err
	}

	req := &model.PatientYearlyBenefitClaimParams{
		PatientID:            pid,
		YearlyBenefitClaimID: yid,
	}

	res, err := c.Usecase.GetByPatientYearlyBenefitClaimID(ctx.Context(), req)
	if err != nil {
		c.Log.WithError(err).Error("GetByPatientYearlyBenefitClaimID failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientYearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient yearly benefit claim retrieved successfully",
		Data:    res,
	})
}

// @Router /patient-yearly-claims [get]
// @Param   page  query int false "Page number" default(1)
// @Param   limit query int false "Items per page" default(10)
// @Success 200 {object} model.PatientYearlyBenefitClaimListResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary List patient yearly benefit claims (paged)
// @Accept json
func (c *PatientYearlyBenefitClaimController) GetAll(ctx *fiber.Ctx) error {
	query := &model.SearchPagingQuery{
		Page:  ctx.QueryInt("page", 0),
		Limit: ctx.QueryInt("limit", 0),
		SearchValue: ctx.Query("search_value", ""),
	}

	items, total, err := c.Usecase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("GetAll patient yearly benefit claims failed")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.PatientYearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient yearly benefit claims fetched successfully",
		Data:    &items,
		Meta:    paging,
	})
}

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/yearly-claims/{yearlyBenefitClaimId} [put]
// @Param  patientId path int true "Patient ID"
// @Param  yearlyBenefitClaimId path int true "Yearly Benefit Claim ID"
// @Success 200 {object} model.PatientYearlyBenefitClaimResponse
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Update patient yearly benefit claim (yearly remaining)
// @Accept json
func (c *PatientYearlyBenefitClaimController) Update(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	ybcIDStr := ctx.Params("yearlyBenefitClaimId")
	if patientIDStr == "" || ybcIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and yearlyBenefitClaimId are required")
	}
	pid, yid, err := c.parseIDs(patientIDStr, ybcIDStr)
	if err != nil {
		return err
	}

	params := &model.PatientYearlyBenefitClaimParams{
		PatientID:            pid,
		YearlyBenefitClaimID: yid,
	}

	res, err := c.Usecase.Update(ctx.Context(), params)
	if err != nil {
		c.Log.WithError(err).Error("Update patient yearly benefit claim failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientYearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient yearly benefit claim updated successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/yearly-claims/{yearlyBenefitClaimId}/reset-remaining [post]
// @Param  patientId path int true "Patient ID"
// @Param  yearlyBenefitClaimId path int true "Yearly Benefit Claim ID"
// @Success 200 {object} model.PatientYearlyBenefitClaimResponse
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Reset yearly remaining by patient & yearly benefit claim
// @Accept json
func (c *PatientYearlyBenefitClaimController) ResetRemainingPlafondByPatientYearlyBenefitClaimID(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	ybcIDStr := ctx.Params("yearlyBenefitClaimId")
	if patientIDStr == "" || ybcIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and yearlyBenefitClaimId are required")
	}
	pid, yid, err := c.parseIDs(patientIDStr, ybcIDStr)
	if err != nil {
		return err
	}

	params := &model.PatientYearlyBenefitClaimParams{
		PatientID:            pid,
		YearlyBenefitClaimID: yid,
	}

	res, err := c.Usecase.ResetRemainingPlafondByPatientYearlyBenefitClaimID(ctx.Context(), params)
	if err != nil {
		c.Log.WithError(err).Error("ResetRemainingPlafondByPatientYearlyBenefitClaimID failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientYearlyBenefitClaimResponse]{
		Code:    fiber.StatusOK,
		Message: "Yearly remaining reset successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-yearly-claims/reset-remaining-all [post]
// @Success 200 {object} model.BaseResponseWrapper
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Reset all yearly remainings (batch)
// @Accept json
func (c *PatientYearlyBenefitClaimController) ResetAllRemainingPlafond(ctx *fiber.Ctx) error {
	if err := c.Usecase.ResetAllRemainingPlafond(ctx.Context()); err != nil {
		c.Log.WithError(err).Error("ResetAllRemainingPlafond failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "All yearly remainings reset successfully",
	})
}

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/reset-remaining [post]
// @Param  patientId path int true "Patient ID"
// @Success 200 {object} model.BaseResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Reset all yearly remainings by patient id (batch)
// @Accept json
func (c *PatientYearlyBenefitClaimController) ResetRemainingPlafondByPatientID(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	if patientIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId is required")
	}
	pid, err := strconv.Atoi(patientIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid value patientId must be number")
	}

	if err := c.Usecase.ResetRemainingPlafondByPatientID(ctx.Context(), uint(pid)); err != nil {
		c.Log.WithError(err).Error("ResetRemainingPlafondByPatientID failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "All yearly remainings by patient id reset successfully",
	})
}

// @Router /api/v1/patient-yearly-claims/patients/{patientId}/yearly-claims/{yearlyBenefitClaimId} [delete]
// @Param  patientId path int true "Patient ID"
// @Param  yearlyBenefitClaimId path int true "Yearly Benefit Claim ID"
// @Success 200 {object} model.BaseResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Yearly Benefit Claims
// @Security BearerAuth
// @Summary Delete patient yearly benefit claim
// @Accept json
func (c *PatientYearlyBenefitClaimController) Delete(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	ybcIDStr := ctx.Params("yearlyBenefitClaimId")
	if patientIDStr == "" || ybcIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and yearlyBenefitClaimId are required")
	}
	pid, yid, err := c.parseIDs(patientIDStr, ybcIDStr)
	if err != nil {
		return err
	}

	params := &model.PatientYearlyBenefitClaimParams{
		PatientID:            uint(pid),
		YearlyBenefitClaimID: uint(yid),
	}

	if err := c.Usecase.Delete(ctx.Context(), params); err != nil {
		c.Log.WithError(err).Error("Delete patient yearly benefit claim failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "Patient yearly benefit claim deleted successfully",
	})
}
