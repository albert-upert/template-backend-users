package postgresql

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/roles/domain"
	"github.com/albert-upert/template-backend-utils-libraries/object"
	"github.com/jackc/pgx/v5"
)

var queryFindAll = `
	SELECT
		id, university_id, name, description, is_active
	FROM iam.roles
	WHERE 1=1
	-- Add dynamic filters here if needed manually or generally where clause
`

// FindAll retrieves a list of roles based on the provided filter.
func (r *Repository) FindAll(ctx context.Context, filter domain.RoleFilter) ([]*domain.Role, int64, error) {
	baseQuery := `
		FROM iam.roles
		WHERE 1=1
	`
	args := pgx.NamedArgs{}

	if filter.UniversityId != "" {
		baseQuery += " AND university_id = @university_id"
		args["university_id"] = filter.UniversityId
	}

	// Simple search implementation
	if filter.Search != "" {
		baseQuery += " AND (name ILIKE @search OR description ILIKE @search)"
		args["search"] = "%" + filter.Search + "%"
	}

	// 1. Count Total
	var total int64
	countQuery := "SELECT count(id)" + baseQuery
	if err := r.db.QueryRow(ctx, countQuery, args).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. Select Data
	selectQuery := `
		SELECT
			id, university_id, name, description, is_active
	` + baseQuery + " ORDER BY id DESC LIMIT @limit OFFSET @offset"

	args["limit"] = filter.Pagination.GetLimit()
	args["offset"] = filter.Pagination.GetOffset()

	rows, err := r.db.Query(ctx, selectQuery, args)
	if err != nil {
		return nil, 0, err
	}

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[RoleEntity])
	if err != nil {
		return nil, 0, err
	}

	roles, err := object.ParseAll[*RoleEntity, *domain.Role](object.TagDB, object.TagObject, records)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}
