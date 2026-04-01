package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
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
