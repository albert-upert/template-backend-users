package usecase

import (
	"context"
	errs "errors"

	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// Get retrieves a user by ID.
func (u *UseCase) Get(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := u.tracer.Start(ctx, "Get")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	user, err := u.repository.FindByID(ctx, id)
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return nil, errors.NotFound("user not found")
		}

		logger.Error().
			Str("func", "repository.FindByID").
			Err(err).
			Msg("failed to find user")

		return nil, errors.DatabaseError("", err)
	}

	return user, nil
}
