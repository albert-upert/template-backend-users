package postgresql

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/redirect/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindUniversityById(ctx context.Context, id string) (*domain.University, error) {
	sql := `
		SELECT
			id, settings
		FROM auth.universities
		WHERE id = $1
	`

	var inst domain.University
	// Using manual scan for safety/simplicity given uncertainties about tags
	if err := r.db.QueryRow(ctx, sql, id).Scan(&inst.Id, &inst.Settings); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("university not found")
		}
		return nil, err
	}

	return &inst, nil
}

func (r *Repository) FindUserBySub(ctx context.Context, universityId string, sub string) (*domain.User, error) {
	sql := `
		SELECT
		    u.id,
		    u.university_id,
		    u.identity_provider,
		    u.external_subject,
		    u.metadata,
		    jsonb_agg(
		        jsonb_build_object(
		            'role_id', r.id,
		            'role_name', r.name,
		            'system_role', r.is_system_role,
		            'institutions', (
		                SELECT array_agg(DISTINCT ur2.university_id)
		                FROM iam.user_roles ur2
		                WHERE ur2.user_id = u.id
		                  AND ur2.role_id = r.id
		            ),
		            'permissions', (
                        SELECT array_agg(DISTINCT p.code)
                        FROM iam.role_permissions rp
                        JOIN iam.permissions p ON p.id = rp.permission_id
                        WHERE rp.role_id = r.id
                    )
		        )
		    ) AS roles
		FROM auth.users u
		JOIN iam.user_roles ur ON u.id = ur.user_id
		JOIN iam.roles r ON ur.role_id = r.id
		WHERE u.external_subject = @subject
		AND u.university_id = @university_id
		AND u.deleted_at IS NULL
		GROUP BY u.id, u.external_subject
		LIMIT 1
	`

	args := pgx.NamedArgs{
		"subject":       sub,
		"university_id": universityId,
	}

	rows, err := r.db.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[domain.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.NotFound("user not found")
		}
		return nil, err
	}

	user.UniversityId = universityId
	return user, nil
}

func (r *Repository) StoreSession(ctx context.Context, session *domain.Session) error {
	sql := `
		INSERT INTO auth.sessions (
			session_id, university_id, user_id, external_subject, roles, access_token, expires_at
		) VALUES (
			@session_id, @university_id, @user_id, @external_subject, @roles, @access_token, @expires_at
		)
	`
	args := pgx.NamedArgs{
		"session_id":       session.SessionId,
		"university_id":    session.UniversityId,
		"user_id":          session.UserId,
		"external_subject": session.ExternalSubject,
		"roles":            session.Roles,
		"access_token":     session.AccessToken,
		"expires_at":       session.ExpiresAt,
	}

	_, err := r.db.Exec(ctx, sql, args)
	return err
}
