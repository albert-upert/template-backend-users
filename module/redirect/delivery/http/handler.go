package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/domain"
)

type RedirectHandler struct {
	useCase domain.RedirectUseCase
}

func NewHandler(useCase domain.RedirectUseCase) *RedirectHandler {
	return &RedirectHandler{
		useCase: useCase,
	}
}

func (h *RedirectHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/redirect/:university_id", h.Redirect)
}

func (h *RedirectHandler) Redirect(c *fiber.Ctx) error {
	ctx := c.UserContext()
	universityId := c.Params("university_id")
	token := c.Cookies("central_access_token", c.Query("access_token"))

	if universityId == "" || token == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail(strconv.Itoa(http.StatusBadRequest), "university_id and token are required"))
	}

	redirectUrl, sessionId, err := h.useCase.Redirect(ctx, universityId, token)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(responses.Fail(strconv.Itoa(http.StatusNotFound), err.Error()))
		}
		if strings.Contains(err.Error(), "invalid token") {
			return c.Status(http.StatusForbidden).JSON(responses.Fail(strconv.Itoa(http.StatusForbidden), err.Error()))
		}

		return c.Status(http.StatusInternalServerError).JSON(responses.Fail(strconv.Itoa(http.StatusInternalServerError), err.Error()))
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect(redirectUrl)
}
