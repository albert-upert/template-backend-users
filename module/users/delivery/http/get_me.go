package http

import (
	"net/http"

	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/gofiber/fiber/v2"
)

type (
	GetMeResponse struct {
		GetUserResponse
		Roles        []middleware.Roles `json:"roles"`
		Institutions []string           `json:"institutions"`
		Permissions  []string           `json:"permissions"`
	}
)

// GetMe handles GET /users/me
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userId := c.Locals(middleware.XUserIdKey).(string)

	user, err := h.useCase.Get(ctx, userId)
	if err != nil {
		return h.handleError(c, err)
	}

	response := GetMeResponse{
		GetUserResponse: GetUserResponse{
			Id:              user.Id,
			UniversityId:    user.UniversityId,
			ExternalSubject: user.ExternalSubject,
			Status:          user.Status,
			Metadata:        user.Metadata,
		},
	}

	if auth, ok := c.Locals(middleware.XUserAuthData).(*middleware.UserRoles); ok {
		response.Roles = auth.Roles
		response.Institutions = auth.Institutions()
		response.Permissions = auth.Permissions()
	}

	return c.Status(http.StatusOK).JSON(responses.Success(response, "User profile retrieved"))
}
