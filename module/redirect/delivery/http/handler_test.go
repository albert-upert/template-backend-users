package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	deliverhttp "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/tests/mocks"
)

func TestRedirectHandler_Redirect(t *testing.T) {
	mockUseCase := new(mocks.RedirectUseCaseMock)
	handler := deliverhttp.NewHandler(mockUseCase)

	app := fiber.New()
	handler.RegisterRoutes(app)

	t.Run("Success", func(t *testing.T) {
		instId := "inst-1"
		token := "valid-token"
		redirectUrl := "http://example.com"
		sessionId := "sess-1"

		mockUseCase.On("Redirect", mock.Anything, instId, token).Return(redirectUrl, sessionId, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+instId+"?token="+token, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.Equal(t, redirectUrl, resp.Header.Get("Location"))

		// Check cookie setup
		cookies := resp.Cookies()
		foundCookie := false
		for _, c := range cookies {
			if c.Name == "session_id" && c.Value == sessionId {
				foundCookie = true
				break
			}
		}
		assert.True(t, foundCookie, "Session cookie should be set")

		mockUseCase.AssertExpectations(t)
	})

	t.Run("MissingParams", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/redirect/inst-1", nil) // Missing token
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		instId := "inst-2"
		token := "token"

		mockUseCase.On("Redirect", mock.Anything, instId, token).Return("", "", errors.New("not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+instId+"?token="+token, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Forbidden", func(t *testing.T) {
		instId := "inst-3"
		token := "invalid"

		mockUseCase.On("Redirect", mock.Anything, instId, token).Return("", "", errors.New("invalid token")).Once()

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+instId+"?token="+token, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		instId := "inst-4"
		token := "token"

		mockUseCase.On("Redirect", mock.Anything, instId, token).Return("", "", errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/redirect/"+instId+"?token="+token, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
