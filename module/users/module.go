package users

import (
	"github.com/albert-upert/template-backend-users/module/users/delivery/http"
	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-users/module/users/repository/postgresql"
	"github.com/albert-upert/template-backend-users/module/users/usecase"
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// Module exports the users module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.UserRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewUserHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.UserHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
