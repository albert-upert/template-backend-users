package roles

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/usecase"
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
