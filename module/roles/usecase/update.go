package usecase

import (
	"context"
	errs "errors"

	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// Update modifies an existing role.
func (u *UseCase) Update(ctx context.Context, role *domain.Role) error {
	ctx, span := u.tracer.Start(ctx, "Update")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Verify exists
	current, err := u.repository.FindByID(ctx, role.Id)
	if err != nil {
		if errs.Is(err, pgx.ErrNoRows) {
			return errors.NotFound("role not found")
		}
		return errors.DatabaseError("", err)
	}

	// Validation: Role names must be unique per university_id (if changed)
	if current.Name != role.Name || current.UniversityId != role.UniversityId {
		existing, err := u.repository.FindByName(ctx, role.UniversityId, role.Name)
		if err != nil && !errs.Is(err, pgx.ErrNoRows) {
			return errors.InternalServerError("failed to validate role name")
		}
		if existing != nil && existing.Id != role.Id {
			return errors.BadRequest("role name already exists in this university")
		}
	}

	if err := u.repository.Update(ctx, role); err != nil {
		logger.Error().
			Str("func", "repository.Update").
			Err(err).
			Msg("failed to update role")

		return errors.DatabaseError("", err)
	}

	return nil
}
