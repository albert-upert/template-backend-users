package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

// MockUseCase is a mock for domain.UseCase
type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) FindAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUseCase) Get(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCase) SyncUser(ctx context.Context, universityId string, token string, code string) (*domain.User, error) {
	args := m.Called(ctx, universityId, token, code)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUseCase) UpdateStatus(ctx context.Context, id string, status string, updatedBy string) error {
	args := m.Called(ctx, id, status, updatedBy)
	return args.Error(0)
}

func (m *MockUseCase) AssignRole(ctx context.Context, cmd domain.AssignRoleCommand) (string, error) {
	args := m.Called(ctx, cmd)
	return args.String(0), args.Error(1)
}

func TestUserHandler_GetMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewUserHandler(mockUseCase, nil) // Auth middleware mocked in test context

		app := fiber.New()

		// Middleware to inject user ID and Auth Data into locals
		app.Use(func(c *fiber.Ctx) error {
			c.Locals(middleware.XUserIdKey, "user-123")
			c.Locals(middleware.XUserAuthData, &middleware.UserRoles{
				UserId: "user-123",
				Roles: []middleware.Roles{
					{
						RoleName:     "admin",
						Institutions: []string{"institution"},
						Permissions:  []string{"users.manage"},
					},
				},
			})
			return c.Next()
		})
		app.Get("/me", handler.GetMe)

		expectedUser := &domain.User{
			Id:              "user-123",
			UniversityId:    "inst-1",
			ExternalSubject: "sub-1",
			Status:          "active",
			Metadata:        map[string]any{"key": "value"},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockUseCase.On("Get", mock.Anything, "user-123").Return(expectedUser, nil)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)

		data := body["data"].(map[string]any)
		assert.Equal(t, "user-123", data["id"])
		assert.Equal(t, "inst-1", data["university_id"])

		// Verify enhanced fields
		roles := data["roles"].([]any)
		assert.Len(t, roles, 1)
		role := roles[0].(map[string]any)
		assert.Equal(t, "admin", role["role_name"])

		groups := data["groups"].([]any)
		assert.Contains(t, groups, "institution")

		permissions := data["permissions"].([]any)
		assert.Contains(t, permissions, "users.manage")
	})

	t.Run("error_user_not_found", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewUserHandler(mockUseCase, nil)

		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals(middleware.XUserIdKey, "user-123")
			return c.Next()
		})
		app.Get("/me", handler.GetMe)

		mockUseCase.On("Get", mock.Anything, "user-123").Return(nil, errors.New("user not found"))

		req := httptest.NewRequest("GET", "/me", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // Or whatever error mapping logic
	})
}
