package usecase

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

// SyncUser performs user synchronization using external IDP.
func (u *UseCase) SyncUser(ctx context.Context, universityId string, token string, code string) (*domain.User, error) {
	ctx, span := u.tracer.Start(ctx, "SyncUser")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// Check Database First
	cachedUser, err := u.repository.FindByExternalSubject(ctx, universityId, code)
	if err == nil && cachedUser != nil {
		return cachedUser, nil
	}
	// Ignore not found error log, proceed to sync

	// Get IDP Client for InstitutionId
	clientIdp, err := u.idp.GetIDP(ctx, universityId)
	if err != nil {
		logger.Error().
			Str("university_id", universityId).
			Err(err).
			Msg("failed to get IDP client")
		return nil, errors.InternalServerError("failed to configure IDP for university")
	}

	// Fetch User from IDP
	logger.Info().Str("code", code).Msg("Fetching user from IDP")

	idpUser, err := clientIdp.GetUserByCode(ctx, token, code)
	if err != nil {
		logger.Error().
			Str("func", "idp.GetUserByCode").
			Err(err).
			Msg("failed to fetch user from IDP")
		return nil, errors.InternalServerError("failed to fetch user from IDP")
	}

	if idpUser == nil {
		return nil, errors.NotFound("user not found in IDP")
	}

	// Convert IDP user struct to map for metadata storage
	metadata := make(map[string]any)
	metaBytes, _ := json.Marshal(idpUser) // Quick marshal/unmarshal to map
	_ = json.Unmarshal(metaBytes, &metadata)

	user := &domain.User{
		UniversityId:     universityId,
		ExternalSubject:  idpUser.Code,
		IdentityProvider: clientIdp.Key(),
		Status:           "active",
		Metadata:         metadata,
	}

	// Upsert User Locally
	if err := u.repository.Store(ctx, user); err != nil {
		logger.Error().Str("func", "repository.Store").Err(err).Msg("failed to store synced user")
		return nil, errors.DatabaseError("", err)
	}

	return user, nil
}
