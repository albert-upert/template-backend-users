package usecase

import (
	"github.com/albert-upert/template-backend-users/module/users/domain"
	"github.com/albert-upert/template-backend-utils-libraries/idp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
