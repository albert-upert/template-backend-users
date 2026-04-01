package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var queryDelete = `DELETE FROM iam.roles WHERE id = @id AND university_id = @university_id`

// Delete removes a role from the database.
func (r *Repository) Delete(ctx context.Context, universityId string, id string) error {
	// First delete permissions
	if err := r.RemovePermissions(ctx, id); err != nil {
		return err
	}

	_, err := r.db.Exec(ctx, queryDelete, pgx.NamedArgs{
		"id":            id,
		"university_id": universityId,
	})

	return err
}
