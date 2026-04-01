package http

import (
	"net/http"

	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateRoleRequest struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		IsActive    bool     `json:"is_active"` // Optional, default true handled logic? TRD doesn't specify default.
		Permissions []string `json:"permissions"`
	}
	CreateRoleResponse struct {
		Id string `json:"id"`
	}
)

// CreateRole handles POST /roles
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	universityId, ok := c.Locals(middleware.XUniversityId).(string)
	if !ok || universityId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing university context"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	// Manual basic validation
	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
	}

	role := domain.Role{
		UniversityId: universityId,
		Name:         req.Name,
		Description:  req.Description,
		IsActive:     req.IsActive,
		Permissions:  req.Permissions,
		CreatedBy:    userId,
		UpdatedBy:    userId,
	}

	if err := h.useCase.Create(ctx, &role); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(CreateRoleResponse{
		Id: role.Id,
	}, "Role created"))
}
