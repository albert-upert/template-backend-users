package postgresql

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/jackc/pgx/v5"
)

var queryStore = `
	INSERT INTO auth.users (
		university_id, external_subject, identity_provider,
		status, metadata, created_at, updated_at
	) VALUES (
		@university_id, @external_subject, @identity_provider,
		@status, @metadata, now(), now()
	)
	ON CONFLICT (university_id, identity_provider, external_subject) WHERE deleted_at IS NULL
	DO UPDATE SET
		metadata = EXCLUDED.metadata,
		updated_at = now()
	RETURNING id
`

// Store upserts a user record.
func (r *Repository) Store(ctx context.Context, user *domain.User) error {
	rows, err := r.db.Query(ctx, queryStore, pgx.NamedArgs{
		"university_id":     user.UniversityId,
		"external_subject":  user.ExternalSubject,
		"identity_provider": user.IdentityProvider,
		"status":            user.Status,
		"metadata":          user.Metadata,
	})
	if err != nil {
		return err
	}

	var id string
	if _, err := pgx.ForEachRow(rows, []any{&id}, func() error { return nil }); err != nil {
		return err
	}
	user.Id = id
	return nil
}
