package usecase

import (
	"context"
	errs "errors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
)

// Create persists a new role record.
func (u *UseCase) Create(ctx context.Context, role *domain.Role) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Validation: Role names must be unique per university_id
	existing, err := u.repository.FindByName(ctx, role.UniversityId, role.Name)
	if err != nil && !errs.Is(err, pgx.ErrNoRows) {
		logger.Error().Err(err).Msg("failed to check role name uniqueness")
		return errors.InternalServerError("failed to validate role")
	}
	if existing != nil {
		return errors.BadRequest("role name already exists in this university")
	}

	if err = u.repository.Store(ctx, role); err != nil {
		logger.Error().
			Str("func", "repository.Store").
			Err(err).
			Msg("failed to store role")

		return errors.DatabaseError("", err)
	}

	return nil
}
