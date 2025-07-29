package http

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type DepartmentController struct {
	DepartmentUseCase *usecase.DepartmentUseCase
	Log *logrus.Logger
}

func NewDepartmentController(usecase *usecase.DepartmentUseCase, log *logrus.Logger) *DepartmentController {
	return &DepartmentController{
		DepartmentUseCase: usecase,
		Log: log,
	}
}

// @Router /api/v1/departments [post]
// @Param  request body model.DepartmentRequest true "Create Department Request"
// @Success 200 {object} model.DepartmentResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Departments
// @Security    BearerAuth api_key
// @Summary Create a new department
// @Description Create a new department with the provided details.
// @Accept json
func (dc *DepartmentController) Create(ctx *fiber.Ctx) error {
	request := &model.DepartmentRequest{}
	ctx.BodyParser(request)

	response, err := dc.DepartmentUseCase.Create(ctx.Context(), request)
	if err != nil {
		dc.Log.WithError(err).Error("Error creating department")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.DepartmentResponse]{
		Code: fiber.StatusCreated,
		Message: "Department created successfully",
		Data: response,
	})
}

// @Router /api/v1/departments/{id} [get]
// @Param  id path int true "Department ID"
// @Success 200 {object} model.DepartmentResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Departments
// @Security    BearerAuth api_key
// @Summary Get a department by ID
// @Description Get a department by its ID.
// @Accept json
func (dc *DepartmentController) GetById(ctx *fiber.Ctx) error {
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

	response, err := dc.DepartmentUseCase.GetById(ctx.Context(), idUint)
	if err != nil {
		dc.Log.WithError(err).Error("Error getting department by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.DepartmentResponse]{
		Code: fiber.StatusOK,
		Message: "Department retrieved successfully",
		Data: response,
	})
}

// @Router /api/v1/departments [get]
// @Success 200 {object} model.DepartmentResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Departments
// @Security    BearerAuth api_key
// @Summary Find departments
// @Description Find departments by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (dc *DepartmentController) GetAll(ctx *fiber.Ctx) error {
	query := &model.PagingQuery{
		Page: ctx.QueryInt("page", 1),
		Limit: ctx.QueryInt("limit", 10),
	}

	responses, total, err := dc.DepartmentUseCase.GetAll(ctx.Context(), query)
	if err != nil {
		dc.Log.WithError(err).Error("Error fetching department types")
		return err
	}

	paging := &model.PaginationPage{
		Page: query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.DepartmentResponse]{
		Code: fiber.StatusOK,
		Message: "Department types fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

// @Router /api/v1/departments/{id} [put]
// @Param  request body model.UpdateDepartmentRequest true "Update Department Request"
// @Param id path string true "Department ID"
// @Success 200 {object} model.DepartmentResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Departments
// @Security    BearerAuth api_key
// @Summary Update a department
// @Description Update a department with the provided details.
// @Accept json
func (dc *DepartmentController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateDepartmentRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		dc.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	request.ID = uint(idInt)

	response, err := dc.DepartmentUseCase.Update(ctx.Context(), request)
	if err != nil {
		dc.Log.WithError(err).Error("Error updating department")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.DepartmentResponse]{
		Code: fiber.StatusOK,
		Message: "Department updated successfully",
		Data: response,
	})
}

// @Router /api/v1/departments/{id} [delete]
// @Param id path string true "Department ID"
// @Success 200 {object} model.DepartmentResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Departments
// @Security    BearerAuth api_key
// @Summary Delete a department
// @Description Delete a department with the provided details.
// @Accept json
func (dc *DepartmentController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		dc.Log.WithError(err).Error("Invalid ID format for deletion")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	err = dc.DepartmentUseCase.Delete(ctx.Context(), uint(idInt))
	if err != nil {
		dc.Log.WithError(err).Error("Error deleting department")
		return err
	}

	return ctx.Status(fiber.StatusNoContent).JSON(model.WebResponse[any]{
		Code: fiber.StatusNoContent,
		Message: "Department deleted successfully",
	})
}