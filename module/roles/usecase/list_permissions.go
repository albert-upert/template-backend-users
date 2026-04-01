package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
)

// ListPermissions retrieves a list of permissions based on filter criteria.
func (u *UseCase) ListPermissions(ctx context.Context, filter domain.PermissionFilter) ([]*domain.Permission, error) {
	ctx, span := u.tracer.Start(ctx, "ListPermissions")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	permissions, err := u.repository.FindAllPermissions(ctx, filter)
	if err != nil {
		logger.Error().
			Str("func", "repository.FindAllPermissions").
			Err(err).
			Msg("failed to find permissions")

		return nil, errors.DatabaseError("", err)
	}

	return permissions, nil
}
