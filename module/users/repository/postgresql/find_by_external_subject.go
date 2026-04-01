package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/object"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

var queryFindBySubject = `
	SELECT
		id, university_id, external_subject, identity_provider,
		status, metadata, created_at, updated_at, deleted_at
	FROM auth.users
	WHERE university_id = @university_id AND external_subject = @subject AND deleted_at IS NULL
	LIMIT 1
`

// FindByExternalSubject retrieves a user by their external subject (sub) within an university.
func (r *Repository) FindByExternalSubject(ctx context.Context, universityId string, subject string) (*domain.User, error) {
	rows, err := r.db.Query(ctx, queryFindBySubject, pgx.NamedArgs{
		"university_id": universityId,
		"subject":       subject,
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
