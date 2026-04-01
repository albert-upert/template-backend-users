package usecase

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/libraries/idp"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/users/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the logic for users module.
type UseCase struct {
	repository domain.UserRepository
	idp        idp.IDPProvider
	tracer     trace.Tracer
}

// NewUseCase creates a new instance of Users UseCase.
func NewUseCase(repository domain.UserRepository, idp idp.IDPProvider) *UseCase {
	return &UseCase{
		repository: repository,
		idp:        idp,
		tracer:     otel.Tracer("users"),
	}
}
