package redirect

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/redirect/usecase"
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
