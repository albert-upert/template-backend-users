package http

import (
	"net/http"

	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/gofiber/fiber/v2"
)

type (
	GetRoleByIDResponse struct {
		Id           string   `json:"id"`
		UniversityId string   `json:"university_id"`
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		IsActive     bool     `json:"is_active"`
		Permissions  []string `json:"permissions"`
	}
)

// GetRoleByID handles GET /roles/:id
func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.BadRequest("Invalid ID"))
	}

	role, err := h.useCase.Get(ctx, id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success(GetRoleByIDResponse{
		Id:           role.Id,
		UniversityId: role.UniversityId,
		Name:         role.Name,
		Description:  role.Description,
		IsActive:     role.IsActive,
		Permissions:  role.Permissions,
	}, "Role retrieved"))
}
