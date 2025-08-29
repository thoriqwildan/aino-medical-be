package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type EmployeeController struct {
	UseCase *usecase.EmployeeUseCase
	Log     *logrus.Logger
}

func NewEmployeeController(useCase *usecase.EmployeeUseCase, log *logrus.Logger) *EmployeeController {
	return &EmployeeController{
		UseCase: useCase,
		Log:     log,
	}
}

// @Router /api/v1/employees [post]
// @Param  request body model.EmployeeRequest true "Create Employee Request"
// @Success 200 {object} model.EmployeeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Employees
// @Security    BearerAuth bearer
// @Summary Create a new employee
// @Description Create a new employee with the provided details.
// @Accept json
func (ec *EmployeeController) Create(ctx *fiber.Ctx) error {
	request := new(model.EmployeeRequest)
	if err := ctx.BodyParser(request); err != nil {
		ec.Log.WithError(err).Error("Error parsing request body in CreateEmployee")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request data")
	}

	response, err := ec.UseCase.Create(ctx.Context(), request)
	if err != nil {
		ec.Log.WithError(err).Error("Error creating employee")
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.EmployeeResponse]{
		Code:    fiber.StatusCreated,
		Message: "Employee created successfully",
		Data:    response,
	})
}

// @Router /api/v1/employees/{id} [get]
// @Param  id path int true "Employee ID"
// @Success 200 {object} model.EmployeeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Employees
// @Security    BearerAuth bearer
// @Summary Get an employee by ID
// @Description Get an employee by its ID.
// @Accept json
func (ec *EmployeeController) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	employeeID, err := strconv.Atoi(id)
	if err != nil {
		ec.Log.WithError(err).Error("Invalid ID format in GetByID")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	response, err := ec.UseCase.GetById(ctx.Context(), uint(employeeID))
	if err != nil {
		ec.Log.WithError(err).Error("Error fetching employee by ID")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.EmployeeResponse]{
		Code:    fiber.StatusOK,
		Message: "Employee fetched successfully",
		Data:    response,
	})
}

// @Router /api/v1/employees [get]
// @Success 200 {object} model.EmployeeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Employees
// @Security    BearerAuth bearer
// @Summary Find employees
// @Description Find employees by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (ec *EmployeeController) GetAll(ctx *fiber.Ctx) error {
	query := &model.SearchPagingQuery{
		SearchValue: ctx.Query("search_value"),
		Page:        ctx.QueryInt("page"),
		Limit:       ctx.QueryInt("limit"),
	}

	responses, total, err := ec.UseCase.GetAll(ctx.Context(), query)
	if err != nil {
		ec.Log.WithError(err).Error("Error fetching employees")
		return err
	}

	paging := &model.PaginationPage{
		Page:  query.Page,
		Limit: query.Limit,
		Total: int(total),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.EmployeeResponse]{
		Code:    fiber.StatusOK,
		Message: "Employees fetched successfully",
		Data:    &responses,
		Meta:    paging,
	})
}

// @Router /api/v1/employees/{id} [put]
// @Param  request body model.UpdateEmployeeRequest true "Update Employee Request"
// @Param id path string true "Employee ID"
// @Success 200 {object} model.EmployeeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Employees
// @Security    BearerAuth bearer
// @Summary Update an employee
// @Description Update an employee with the provided details.
// @Accept json
func (c *EmployeeController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateEmployeeRequest)
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.EmployeeResponse]{
		Code:    fiber.StatusOK,
		Message: "Employee updated successfully",
		Data:    response,
	})
}

// @Router /api/v1/employees/{id} [delete]
// @Param id path string true "Employee ID"
// @Success 200 {object} model.EmployeeResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Employees
// @Security    BearerAuth bearer
// @Summary Delete an employee
// @Description Delete an employee with the provided details.
// @Accept json
func (c *EmployeeController) Delete(ctx *fiber.Ctx) error {
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[any]{
		Code:    fiber.StatusOK,
		Message: "Employee deleted successfully",
	})
}
