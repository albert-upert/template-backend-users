package mocks

import (
	"context"

	"github.com/albert-upert/template-backend-users/module/redirect/domain"
	"github.com/stretchr/testify/mock"
)

// RedirectUseCaseMock is a mock for RedirectUseCase
type RedirectUseCaseMock struct {
	mock.Mock
}

func (m *RedirectUseCaseMock) Redirect(ctx context.Context, universityId string, token string) (string, string, error) {
	args := m.Called(ctx, universityId, token)
	return args.String(0), args.String(1), args.Error(2)
}

// RedirectRepositoryMock is a mock for RedirectRepository
type RedirectRepositoryMock struct {
	mock.Mock
}

func (m *RedirectRepositoryMock) FindUniversityById(ctx context.Context, id string) (*domain.University, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.University), args.Error(1)
}

func (m *RedirectRepositoryMock) FindUserBySub(ctx context.Context, universityId string, sub string) (*domain.User, error) {
	args := m.Called(ctx, universityId, sub)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *RedirectRepositoryMock) StoreSession(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}
