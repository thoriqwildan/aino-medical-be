package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoriqwildan/aino-medical-be/internal/model"
	"github.com/thoriqwildan/aino-medical-be/internal/usecase"
)

type FamilyMemberController struct {
	UseCase *usecase.FamilyMemberUseCase
	Log *logrus.Logger
	Config *viper.Viper
}

func NewFamilyMemberController(useCase *usecase.FamilyMemberUseCase, log *logrus.Logger, config *viper.Viper) *FamilyMemberController {
	return &FamilyMemberController{
		UseCase: useCase,
		Log: log,
		Config: config,
	}
}

// @Router /api/v1/family-members [post]
// @Param  request body model.FamilyMemberRequest true "Create Family Member Request"
// @Success 200 {object} model.FamilyMemberResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Family Members
// @Security    BearerAuth api_key
// @Summary Create a new family member
// @Description Create a new family member with the provided details.
// @Accept json
func (c *FamilyMemberController) Create(ctx *fiber.Ctx) error {
	request := new(model.FamilyMemberRequest)

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("Failed to parse request body")
		return err
	}

	familyMember, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Failed to create family member")
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[model.FamilyMemberResponse]{
		Code: fiber.StatusCreated,
		Message: "Family member created successfully",
		Data: familyMember,
	})
}

// @Router /api/v1/family-members/{id} [get]
// @Param  id path int true "Family Member ID"
// @Success 200 {object} model.FamilyMemberResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Family Members
// @Security    BearerAuth api_key
// @Summary Get a family member by ID
// @Description Get a family member by its ID.
// @Accept json
func (c *FamilyMemberController) GetById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}

	familyMemberID, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format in GetByID")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	familyMember, err := c.UseCase.GetByID(ctx.Context(), uint(familyMemberID))
	if err != nil {
		c.Log.WithError(err).Error("Failed to get family member by ID")
		return err
	}

	return ctx.JSON(model.WebResponse[model.FamilyMemberResponse]{
		Code: fiber.StatusOK,
		Message: "Family member retrieved successfully",
		Data: familyMember,
	})
}

// @Router /api/v1/family-members [get]
// @Success 200 {object} model.FamilyMemberResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Family Members
// @Security    BearerAuth api_key
// @Summary Find family members
// @Description Find family members by their attributes.
// @Param   page query     int               false       "Page number" default(1)
// @Param   limit query    int               false       "Number of items per page" default
// @Accept json
func (c *FamilyMemberController) GetAll(ctx *fiber.Ctx) error {
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

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[[]model.FamilyMemberResponse]{
		Code: fiber.StatusOK,
		Message: "Family members fetched successfully",
		Data: &responses,
		Meta: paging,
	})
}

// @Router /api/v1/family-members/{id} [put]
// @Param  request body model.UpdateFamilyMemberRequest true "Update Family Member Request"
// @Param id path string true "Family Member ID"
// @Success 200 {object} model.FamilyMemberResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Family Members
// @Security    BearerAuth api_key
// @Summary Update a family member
// @Description Update a family member with the provided details.
// @Accept json
func (c *FamilyMemberController) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "ID is required")
	}

	request := new(model.UpdateFamilyMemberRequest)
	ctx.BodyParser(&request)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Log.WithError(err).Error("Invalid ID format for update")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}
	idUint := uint(idInt)
	request.ID = &idUint

	response, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		c.Log.WithError(err).Error("Error updating family member")
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[model.FamilyMemberResponse]{
		Code: fiber.StatusOK,
		Message: "Family member updated successfully",
		Data: response,
	})
}

// @Router /api/v1/family-members/{id} [delete]
// @Param id path string true "Family Member ID"
// @Success 200 {object} model.FamilyMemberResponseWrapper
// @Failure 400 {object} model.ErrorWrapper "Bad Request"
// @Failure 500 {object} model.ErrorWrapper "Internal Server Error"
// @Tags Family Members
// @Security    BearerAuth api_key
// @Summary Delete a family member
// @Description Delete a family member with the provided details.
// @Accept json
func (c *FamilyMemberController) Delete(ctx *fiber.Ctx) error {
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