package postgresql

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-utils-libraries/object"
	"github.com/jackc/pgx/v5"
)

var queryFindById = `
	SELECT
		id, university_id, name, description, is_active
	FROM iam.roles
	WHERE id = @id
	LIMIT 1
`

// FindByID retrieves a single role by their ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[RoleEntity])
	if err != nil {
		return nil, err
	}

	role, err := object.Parse[*RoleEntity, *domain.Role](object.TagDB, object.TagObject, record)
	if err != nil {
		return nil, err
	}

	// Fetch permissions
	perms, err := r.GetPermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Permissions = perms

	return role, nil
}

var queryFindByName = `
	SELECT
		id, university_id, name, description, is_active
	FROM iam.roles
	WHERE university_id = @university_id AND name = @name
	LIMIT 1
`

// FindByName retrieves a single role by name and university.
func (r *Repository) FindByName(ctx context.Context, universityId string, name string) (*domain.Role, error) {
	rows, err := r.db.Query(ctx, queryFindByName, pgx.NamedArgs{
		"university_id": universityId,
		"name":          name,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[RoleEntity])
	if err != nil {
		return nil, err
	}

	role, err := object.Parse[*RoleEntity, *domain.Role](object.TagDB, object.TagObject, record)
	if err != nil {
		return nil, err
	}

	return role, nil
}
