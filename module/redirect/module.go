package redirect

import (
	"github.com/albert-upert/template-backend-users/module/redirect/delivery/http"
	"github.com/albert-upert/template-backend-users/module/redirect/domain"
	"github.com/albert-upert/template-backend-users/module/redirect/repository/postgresql"
	"github.com/albert-upert/template-backend-users/module/redirect/usecase"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"redirect",
	fx.Provide(
		fx.Annotate(
			postgresql.NewRepository,
			fx.As(new(domain.RedirectRepository)),
		),
		fx.Annotate(
			usecase.NewUseCase,
			fx.As(new(domain.RedirectUseCase)),
		),
		http.NewHandler,
	),
	fx.Invoke(
		registerRoutes,
	),
)

func registerRoutes(h *http.RedirectHandler, app *fiber.App) {
	h.RegisterRoutes(app)
}
