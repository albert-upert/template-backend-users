package http

import (
	"net/http"

	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/albert-upert/template-backend-utils-libraries/middleware"
	"github.com/albert-upert/template-backend-utils-libraries/responses"
	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests for users module.
type UserHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the users module.
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/users", middleware.TraceMiddleware)

	group.Get("/me", h.auth.Authenticate(), h.GetMe)
	group.Get("/", h.auth.Authenticate("users.manage.all.view"), h.GetUsers)

	group.Post("/", h.auth.Authenticate("users.manage.all.edit"), h.SyncUser)

	group.Patch("/:id/status", h.auth.Authenticate("users.manage.all.edit"), h.UpdateStatus)
	group.Post("/:id/roles", h.auth.Authenticate("users.manage.all.edit"), h.AssignRole)
}

// handleError handles errors by mapping them to standardized responses.
func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
