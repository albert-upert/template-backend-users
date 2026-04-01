package http

import (
	"net/http"

	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/gofiber/fiber/v2"
)

type (
	SyncUserRequest struct {
		Code string `json:"code" validate:"required"`
	}
	SyncUserResponse struct {
		Id string `json:"id"`
	}
)

// SyncUser handles POST /users
func (h *UserHandler) SyncUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Get university_id from locals
	universityId, ok := c.Locals(middleware.XUniversityId).(string)
	if !ok || universityId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing university context"))
	}

	var req SyncUserRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	if req.Code == "" {
		return h.handleError(c, errors.BadRequest("code is required"))
	}

	// Get token from locals
	token, ok := c.Locals(middleware.XTokenKey).(string)
	if !ok || token == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing token"))
	}

	user, err := h.useCase.SyncUser(ctx, universityId, token, req.Code)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(SyncUserResponse{
		Id: user.Id,
	}, "User synced"))
}
