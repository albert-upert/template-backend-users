package domain

import (
	"context"
)

// RedirectUseCase defines the business logic contract for Redirect module.
type RedirectUseCase interface {
	Redirect(ctx context.Context, universityId string, token string) (string, string, error) // Returns redirectUrl, sessionId, error
}
