package users

import (
	gofiber "github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/repository/postgresql"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/usecase"
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
