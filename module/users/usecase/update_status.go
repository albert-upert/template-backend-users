package usecase

import (
	"context"

	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/rs/zerolog"
)

// UpdateStatus updates a user's status.
func (u *UseCase) UpdateStatus(ctx context.Context, id string, status string, updatedBy string) error {
	ctx, span := u.tracer.Start(ctx, "UpdateStatus")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Validate status enum
	validStatuses := map[string]bool{"active": true, "inactive": true, "suspended": true, "pending": true}
	if !validStatuses[status] {
		return errors.BadRequest("invalid status")
	}

	if err := u.repository.UpdateStatus(ctx, id, status, updatedBy); err != nil {
		logger.Error().
			Str("func", "repository.UpdateStatus").
			Err(err).
			Msg("failed to update user status")
		return errors.DatabaseError("", err)
	}

	return nil
}
