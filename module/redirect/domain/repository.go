package domain

import (
	"context"
	"time"

	"github.com/albert-upert/template-backend-utils-libraries/idp"
)

// University represents the university settings for redirect.
type University struct {
	Id          string      `object:"id"`
	RedirectUrl string      `object:"redirect_url"`
	Settings    idp.Setting `object:"settings"`
}

// Session represents an authenticated user session.
type Session struct {
	SessionId       string    `object:"session_id"`
	UniversityId    string    `object:"university_id"`
	UserId          string    `object:"user_id"`
	ExternalSubject string    `object:"external_subject"`
	Roles           any       `object:"roles"` // Changed to any to support JSONB structure
	AccessToken     string    `object:"access_token"`
	ExpiresAt       time.Time `object:"expires_at"`
}

// User represents the user found by sub.
type User struct {
	Id               string `object:"id"`
	UniversityId     string `object:"university_id"`
	ExternalSubject  string `object:"external_subject"`
	IdentityProvider string `object:"identity_provider"`
	Metadata         any    `object:"metadata"`
	Roles            any    `object:"roles"`
}

// RedirectRepository defines the persistence layer contract.
type RedirectRepository interface {
	FindUniversityById(ctx context.Context, id string) (*University, error)
	FindUserBySub(ctx context.Context, universityId string, sub string) (*User, error)
	StoreSession(ctx context.Context, session *Session) error
}
