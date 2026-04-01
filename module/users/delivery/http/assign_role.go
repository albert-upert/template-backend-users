package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

type (
	AssignRoleRequest struct {
		UniversityId  string `json:"university_id" validate:"required"`
		RoleId        string `json:"role_id" validate:"required"`
		InstitutionId string `json:"institution_id" validate:"required"`
	}
	AssignRoleResponse struct {
		Id string `json:"id"`
	}
)

// AssignRole handles POST /users/:id/roles
func (h *UserHandler) AssignRole(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	userId, ok := c.Locals(middleware.XUserIdKey).(string)
	if !ok || userId == "" {
		return h.handleError(c, errors.Unauthorized("Invalid or missing user context"))
	}

	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.BadRequest("Invalid request body"))
	}

	cmd := domain.AssignRoleCommand{
		UserId:        id,
		RoleId:        req.RoleId,
		UniversityId:  req.UniversityId,
		InstitutionId: req.InstitutionId,
		AssignedBy:    userId,
	}

	assignmentId, err := h.useCase.AssignRole(ctx, cmd)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(AssignRoleResponse{
		Id: assignmentId,
	}, "Role assigned"))
}
