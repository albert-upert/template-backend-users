package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
)

// Delete removes a role from the system.
func (u *UseCase) Delete(ctx context.Context, universityId string, id string) error {
	ctx, span := u.tracer.Start(ctx, "Delete")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	if err := u.repository.Delete(ctx, universityId, id); err != nil {
		logger.Error().
			Str("func", "repository.Delete").
			Err(err).
			Msg("failed to delete role")

		return errors.DatabaseError("", err)
	}

	return nil
}
