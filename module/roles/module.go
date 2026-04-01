package roles

import (
	"github.com/albert-upert/template-backend-users/module/roles/delivery/http"
	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-users/module/roles/repository/postgresql"
	"github.com/albert-upert/template-backend-users/module/roles/usecase"
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// Module exports the roles module for Fx.
var Module = fx.Options(
	fx.Provide(
		postgresql.NewRepository,
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.RoleRepository)),
		),
		usecase.NewUseCase,
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.UseCase)),
		),
		http.NewRoleHandler,
	),
	fx.Invoke(registerRoutes),
)

func registerRoutes(h *http.RoleHandler, app *gofiber.App) {
	h.RegisterRoutes(app)
}
