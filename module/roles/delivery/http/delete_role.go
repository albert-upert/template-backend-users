package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/responses"
)

// DeleteRole handles DELETE /roles/:id
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	universityId, _ := c.Locals(middleware.XUniversityId).(string)

	if err := h.useCase.Delete(ctx, universityId, id); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Role deleted"))
}
