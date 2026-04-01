package usecase

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/rs/zerolog"
)

// FindAll retrieves a list of users based on filter criteria.
func (u *UseCase) FindAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int64, error) {
	ctx, span := u.tracer.Start(ctx, "FindAll")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	users, total, err := u.repository.FindAll(ctx, filter)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindAll").
			Err(err).
			Msg("failed to find users")

		return nil, 0, errors.DatabaseError("", err)
	}

	return users, total, nil
}
