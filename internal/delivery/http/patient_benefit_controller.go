package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type PatientBenefitController struct {
	Log     *logrus.Logger
	Usecase *usecase.PatientBenefitUsecase
}

func NewPatientBenefitController(uc *usecase.PatientBenefitUsecase, log *logrus.Logger) *PatientBenefitController {
	return &PatientBenefitController{
		Log:     log,
		Usecase: uc,
	}
}

func (c *PatientBenefitController) parseIDs(patientIDStr, benefitIDStr string) (uint, uint, error) {
	pid, err := strconv.Atoi(patientIDStr)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusBadRequest, "invalid patientId format")
	}
	bid, err := strconv.Atoi(benefitIDStr)
	if err != nil {
		return 0, 0, fiber.NewError(fiber.StatusBadRequest, "invalid benefitId format")
	}
	return uint(pid), uint(bid), nil
}

// @Router /api/v1/patient-benefits [post]
// @Param  request body model.PatientBenefitParams true "Create Patient Benefit"
// @Success 201 {object} model.PatientBenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Create patient benefit (find or create for current year)
// @Accept json
func (c *PatientBenefitController) Create(ctx *fiber.Ctx) error {
	req := new(model.PatientBenefitParams)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	res, err := c.Usecase.Create(ctx.Context(), req)
	if err != nil {
		c.Log.WithError(err).Error("Create patient benefit failed")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.PatientBenefitResponse]{
		Code:    fiber.StatusCreated,
		Message: "Patient benefit created successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-benefits/patients/{patientId}/benefits/{benefitId} [get]
// @Param  patientId path int true "Patient ID"
// @Param  benefitId path int true "Benefit ID"
// @Success 200 {object} model.PatientBenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Get patient benefit by patientId & benefitId
// @Accept json
func (c *PatientBenefitController) GetByPatientBenefitID(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	benefitIDStr := ctx.Params("benefitId")
	if patientIDStr == "" || benefitIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and benefitId are required")
	}

	pid, bid, err := c.parseIDs(patientIDStr, benefitIDStr)
	if err != nil {
		return err
	}

	req := &model.PatientBenefitParams{
		PatientID: pid,
		BenefitID: bid,
	}

	res, err := c.Usecase.GetByPatientBenefitID(ctx.Context(), req)
	if err != nil {
		c.Log.WithError(err).Error("GetByPatientBenefitID failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientBenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient benefit retrieved successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-benefits [get]
// @Param   page  query int false "Page number" default(1)
// @Param   limit query int false "Items per page" default(10)
// @Success 200 {object} model.PatientBenefitListResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary List patient benefits (paged)
// @Description Returns paginated patient benefits for the current filters.
// @Accept json
func (c *PatientBenefitController) GetAll(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page:  ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	items, total, err := c.Usecase.GetAll(ctx.Context(), query)
	if err != nil {
		c.Log.WithError(err).Error("GetAll patient benefits failed")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]*model.PatientBenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient benefits fetched successfully",
		Data:    &items,
		Meta:    paging,
	})
}

// @Router /api/v1/patient-benefits/patients/{patientId}/benefits/{benefitId} [put]
// @Param  patientId path int true "Patient ID"
// @Param  benefitId path int true "Benefit ID"
// @Param  request body model.UpdatePatientBenefitRequest true "Update Patient Benefit"
// @Success 200 {object} model.PatientBenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Update patient benefit
// @Accept json
func (c *PatientBenefitController) Update(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	benefitIDStr := ctx.Params("benefitId")
	if patientIDStr == "" || benefitIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and benefitId are required")
	}
	pid, bid, err := c.parseIDs(patientIDStr, benefitIDStr)
	if err != nil {
		return err
	}

	body := new(model.UpdatePatientBenefitRequest)
	if err := ctx.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	params := &model.PatientBenefitParams{
		PatientID: pid,
		BenefitID: bid,
	}

	res, err := c.Usecase.Update(ctx.Context(), params, body)
	if err != nil {
		c.Log.WithError(err).Error("Update patient benefit failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientBenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient benefit updated successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-benefits/patients/{patientId}/benefits/{benefitId}/reset-remaining [post]
// @Param  patientId path int true "Patient ID"
// @Param  benefitId path int true "Benefit ID"
// @Success 200 {object} model.PatientBenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Reset remaining plafond for a specific patient benefit
// @Accept json
func (c *PatientBenefitController) ResetRemainingPlafondByPatientBenefitID(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	benefitIDStr := ctx.Params("benefitId")
	if patientIDStr == "" || benefitIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and benefitId are required")
	}
	pid, bid, err := c.parseIDs(patientIDStr, benefitIDStr)
	if err != nil {
		return err
	}

	params := &model.PatientBenefitParams{
		PatientID: pid,
		BenefitID: bid,
	}

	res, err := c.Usecase.ResetRemainingPlafondByPatientBenefitID(ctx.Context(), params)
	if err != nil {
		c.Log.WithError(err).Error("ResetRemainingPlafondByPatientBenefitID failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientBenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Remaining plafond reset successfully",
		Data:    res,
	})
}

// @Router /api/v1/patient-benefits/reset-all [post]
// @Success 200 {object} model.BaseResponseWrapper
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Reset all remaining plafonds (batch)
// @Accept json
func (c *PatientBenefitController) ResetAllRemainingPlafond(ctx *fiber.Ctx) error {
	if err := c.Usecase.ResetAllRemainingPlafond(ctx.Context()); err != nil {
		c.Log.WithError(err).Error("ResetAllRemainingPlafond failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "All remaining plafonds reset successfully",
	})
}

// @Router /api/v1/patient-benefits/patients/{patientId}/benefits/{benefitId} [delete]
// @Param  patientId path int true "Patient ID"
// @Param  benefitId path int true "Benefit ID"
// @Success 200 {object} model.PatientBenefitResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 404 {object} model.ErrorWrapper "Not Found"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Patient Benefits
// @Security BearerAuth
// @Summary Delete patient benefit
// @Accept json
func (c *PatientBenefitController) Delete(ctx *fiber.Ctx) error {
	patientIDStr := ctx.Params("patientId")
	benefitIDStr := ctx.Params("benefitId")
	if patientIDStr == "" || benefitIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "patientId and benefitId are required")
	}
	pid, bid, err := c.parseIDs(patientIDStr, benefitIDStr)
	if err != nil {
		return err
	}

	params := &model.PatientBenefitParams{
		PatientID: pid,
		BenefitID: bid,
	}

	res, err := c.Usecase.Delete(ctx.Context(), params)
	if err != nil {
		c.Log.WithError(err).Error("Delete patient benefit failed")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.PatientBenefitResponse]{
		Code:    fiber.StatusOK,
		Message: "Patient benefit deleted successfully",
		Data:    res,
	})
}
