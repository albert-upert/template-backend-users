package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

var queryAssignRole = `
	INSERT INTO iam.user_roles (
		university_id, user_id, role_id, university_id,
		assigned_at, is_active, assigned_by
	) VALUES (
		@university_id, @user_id, @role_id, @university_id,
		now(), true, @assigned_by
	)
	ON CONFLICT (university_id, user_id, role_id, university_id)
	DO UPDATE SET
		is_active = true,
		assigned_at = now(), -- Re-activate if exists
		assigned_by = EXCLUDED.assigned_by
	RETURNING id
`

// AssignRole creates a new user-role assignment.
func (r *Repository) AssignRole(ctx context.Context, role *domain.UserRole) error {
	rows, err := r.db.Query(ctx, queryAssignRole, pgx.NamedArgs{
		"university_id":  role.UniversityId,
		"user_id":        role.UserId,
		"role_id":        role.RoleId,
		"institution_id": role.InstitutionId,
		"assigned_by":    role.AssignedBy,
	})
	if err != nil {
		return err
	}

	var id string
	if _, err := pgx.ForEachRow(rows, []any{&id}, func() error { return nil }); err != nil {
		return err
	}
	role.Id = id
	return nil
}
