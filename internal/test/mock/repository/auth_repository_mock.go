package mock

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth/dto"
	"be-dashboard-nba/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func NewAuthRepository() *AuthRepositoryMock { return &AuthRepositoryMock{} }

func (m *AuthRepositoryMock) CreateSessionQuery(ctx context.Context, args dto.SessionParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) UpdateSessionQuery(ctx context.Context, args dto.SessionParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) DeleteSessionQuery(ctx context.Context, id string) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}

func (m *AuthRepositoryMock) ReadDetailSessionQuery(ctx context.Context, id string) (model.Session, error) {
	call := m.Called(ctx, id)

	var r0 model.Session
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Session); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *AuthRepositoryMock) ReadDetailUserByEmailQuery(ctx context.Context, email string) (model.User, error) {
	call := m.Called(ctx, email)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *AuthRepositoryMock) CreateLoginAttemp(ctx context.Context, args dto.LoginAttempParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) CreateLoginRecord(ctx context.Context, args dto.LoginRecord) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *AuthRepositoryMock) CheckPermissionQuery(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error) {
	call := m.Called(ctx, menuURL, userID, permissionCode)
	return call.Bool(0), call.Error(1)
}

func (m *AuthRepositoryMock) ReadDetailUserByIdQuery(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}
