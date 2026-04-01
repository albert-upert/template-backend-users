package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
)

type (
	UpdateRoleRequest struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		IsActive    bool     `json:"is_active"`
		Permissions []string `json:"permissions"`
	}
)

// UpdateRole handles PUT /roles/:id
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	universityId, ok := c.Locals(middleware.XUniversityId).(string)
	if !ok || universityId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing university context"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if req.Name == "" {
		return h.handleError(c, errors.BadRequest("field name is required"))
	}

	role := domain.Role{
		Id:           id,
		UniversityId: universityId,
		Name:         req.Name,
		Description:  req.Description,
		IsActive:     req.IsActive,
		Permissions:  req.Permissions,
		UpdatedBy:    userId,
	}

	if err := h.useCase.Update(ctx, &role); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Role updated"))
}
