package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/responses"
)

type (
	UpdateStatusRequest struct {
		Status string `json:"status" validate:"required"`
	}
)

// UpdateStatus handles PATCH /users/:id/status
func (h *UserHandler) UpdateStatus(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if err := h.useCase.UpdateStatus(ctx, id, req.Status, userId); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "User status updated"))
}
