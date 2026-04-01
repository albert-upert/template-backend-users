package usecase

import (
	"context"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

// AssignRole assigns a role to a user.
func (u *UseCase) AssignRole(ctx context.Context, cmd domain.AssignRoleCommand) (string, error) {
	ctx, span := u.tracer.Start(ctx, "AssignRole")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Validate inputs
	if cmd.UserId == "" || cmd.RoleId == "" || cmd.UniversityId == "" {
		return "", errors.BadRequest("missing required fields")
	}

	if cmd.UniversityId == "" {
		return "", errors.BadRequest("university_id is currently required")
	}

	userRole := &domain.UserRole{
		UserId:        cmd.UserId,
		RoleId:        cmd.RoleId,
		UniversityId:  cmd.UniversityId,
		InstitutionId: cmd.InstitutionId,
		AssignedBy:    &cmd.AssignedBy,
		IsActive:      true,
	}

	if err := u.repository.AssignRole(ctx, userRole); err != nil {
		logger.Error().
			Str("func", "repository.AssignRole").
			Err(err).
			Msg("failed to assign role")
		return "", errors.DatabaseError("", err)
	}

	return userRole.Id, nil
}
