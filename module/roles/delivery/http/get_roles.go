package http

import (
	"net/http"
	"strconv"

	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/albert-upert/template-backend-utils-libraries/types"
	"github.com/gofiber/fiber/v2"
)

type (
	GetRolesResponse struct {
		Id           string `json:"id"`
		UniversityId string `json:"university_id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		IsActive     bool   `json:"is_active"`
	}
)

// GetRoles handles GET /roles
func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	ctx := c.UserContext()

	universityId, _ := c.Locals(middleware.XUniversityId).(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	filter := domain.RoleFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: pageSize,
		},
		UniversityId: universityId,
		Search:       c.Query("search"),
	}

	roles, total, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]GetRolesResponse, len(roles))
	for i, r := range roles {
		result[i] = GetRolesResponse{
			Id:           r.Id,
			UniversityId: r.UniversityId,
			Name:         r.Name,
			Description:  r.Description,
			IsActive:     r.IsActive,
		}
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       pageSize,
		Total:      total,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessWithMeta(result, "Roles retrieved", meta))
}
