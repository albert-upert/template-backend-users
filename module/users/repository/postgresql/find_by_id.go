package postgresql

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-utils-libraries/object"
	"github.com/jackc/pgx/v5"
)

var queryFindById = `
	SELECT
		id, university_id, external_subject, identity_provider,
		status, metadata, created_at, updated_at, deleted_at
	FROM auth.users
	WHERE id = @id AND deleted_at IS NULL
	LIMIT 1
`

// FindByID retrieves a user by ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[UserEntity])
	if err != nil {
		return nil, err
	}

	return object.Parse[*UserEntity, *domain.User](object.TagDB, object.TagObject, record)
}
