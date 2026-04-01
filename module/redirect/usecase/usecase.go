package usecase

import (
	"context"
	"time"

	"github.com/albert-upert/template-backend-users/config"
	"github.com/albert-upert/template-backend-users/module/redirect/domain"
	"github.com/albert-upert/template-backend-utils-libraries/errors"
	"github.com/albert-upert/template-backend-utils-libraries/idp"
	"github.com/albert-upert/template-backend-utils-libraries/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var _ domain.RedirectUseCase = (*UseCase)(nil)

type UseCase struct {
	defaultRedirectUrl string
	repository         domain.RedirectRepository
	idp                idp.IDPProvider
	tracer             trace.Tracer
}

func NewUseCase(app *config.InternalAppConfig, repository domain.RedirectRepository, idp idp.IDPProvider) *UseCase {
	return &UseCase{
		defaultRedirectUrl: app.RedirectUrl,
		repository:         repository,
		idp:                idp,
		tracer:             otel.Tracer("redirect"),
	}
}

func (u *UseCase) Redirect(ctx context.Context, universityId string, token string) (string, string, error) {
	ctx, span := u.tracer.Start(ctx, "Redirect")
	defer span.End()

	university, err := u.repository.FindUniversityById(ctx, universityId)
	if err != nil {
		return "", "", err
	}

	idpClient, err := u.idp.GetIDP(ctx, universityId)
	if err != nil {
		return "", "", errors.InternalServerError("failed to get idp client")
	}

	authSession, err := idpClient.Check(ctx, token)
	if err != nil {
		return "", "", errors.New(errors.ErrorTypeUnauthorized, 403, "invalid token", nil)
	}

	user, err := u.repository.FindUserBySub(ctx, universityId, authSession.Sub)
	if err != nil {
		return "", "", err
	}

	// Calculate session expiry based on authSession.ExpiresIn (int seconds)
	expiresAt := time.Now().Add(time.Duration(authSession.ExpiresIn) * time.Second)

	session := &domain.Session{
		SessionId:       types.GenerateID(),
		UniversityId:    universityId,
		UserId:          user.Id,
		ExternalSubject: user.ExternalSubject,
		Roles:           user.Roles,
		AccessToken:     token,
		ExpiresAt:       expiresAt,
	}

	if err := u.repository.StoreSession(ctx, session); err != nil {
		return "", "", errors.DatabaseError("", err)
	}

	redirectUrl := university.Settings.Url
	if redirectUrl == "" {
		redirectUrl = u.defaultRedirectUrl
	}

	return redirectUrl, session.SessionId, nil
}
