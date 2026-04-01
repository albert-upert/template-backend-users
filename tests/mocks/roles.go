package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/user/module/roles/domain"
)

// RolesUseCaseMock is a mock for Roles UseCase
type RolesUseCaseMock struct {
	mock.Mock
}

func (m *RolesUseCaseMock) FindAll(ctx context.Context, filter domain.RoleFilter) ([]*domain.Role, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Role), args.Get(1).(int64), args.Error(2)
}

func (m *RolesUseCaseMock) Get(ctx context.Context, id string) (*domain.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *RolesUseCaseMock) Create(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *RolesUseCaseMock) Update(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *RolesUseCaseMock) Delete(ctx context.Context, universityId string, id string) error {
	args := m.Called(ctx, universityId, id)
	return args.Error(0)
}

func (m *RolesUseCaseMock) ListPermissions(ctx context.Context, filter domain.PermissionFilter) ([]*domain.Permission, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}

// RolesRepositoryMock is a mock for Roles Repository
type RolesRepositoryMock struct {
	mock.Mock
}

func (m *RolesRepositoryMock) FindAll(ctx context.Context, filter domain.RoleFilter) ([]*domain.Role, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Role), args.Get(1).(int64), args.Error(2)
}

func (m *RolesRepositoryMock) FindByID(ctx context.Context, id string) (*domain.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *RolesRepositoryMock) FindByName(ctx context.Context, universityId string, name string) (*domain.Role, error) {
	args := m.Called(ctx, universityId, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *RolesRepositoryMock) Store(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *RolesRepositoryMock) Update(ctx context.Context, role *domain.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *RolesRepositoryMock) Delete(ctx context.Context, universityId string, id string) error {
	args := m.Called(ctx, universityId, id)
	return args.Error(0)
}

func (m *RolesRepositoryMock) AddPermissions(ctx context.Context, roleId string, permissions []string) error {
	args := m.Called(ctx, roleId, permissions)
	return args.Error(0)
}

func (m *RolesRepositoryMock) RemovePermissions(ctx context.Context, roleId string) error {
	args := m.Called(ctx, roleId)
	return args.Error(0)
}

func (m *RolesRepositoryMock) GetPermissions(ctx context.Context, roleId string) ([]string, error) {
	args := m.Called(ctx, roleId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *RolesRepositoryMock) FindAllPermissions(ctx context.Context, filter domain.PermissionFilter) ([]*domain.Permission, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Permission), args.Error(1)
}
