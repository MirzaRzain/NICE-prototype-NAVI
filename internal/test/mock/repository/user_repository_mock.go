package mock

import (
	"context"

	"be-dashboard-nba/internal/application/user/dto"
	"be-dashboard-nba/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepository() *UserRepositoryMock { return &UserRepositoryMock{} }

func (m *UserRepositoryMock) CreateUserQuery(ctx context.Context, args dto.CreateUserParams) (string, error) {
	call := m.Called(ctx, args)

	var r0 string
	raw0 := call.Get(0)
	if casted, ok := raw0.(string); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) CreateUserRoleQuery(ctx context.Context, roleID int, userID string) error {
	call := m.Called(ctx, roleID, userID)
	return call.Error(0)
}

func (m *UserRepositoryMock) UpdateUserQuery(ctx context.Context, args dto.UpdateUserParams) error {
	call := m.Called(ctx, args)
	return call.Error(0)
}

func (m *UserRepositoryMock) UpdateUserRoleQuery(ctx context.Context, roleID int, userID string) error {
	call := m.Called(ctx, roleID, userID)
	return call.Error(0)
}

func (m *UserRepositoryMock) DeleteUserQuery(ctx context.Context, id string) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}

func (m *UserRepositoryMock) ReadUsersQuery(ctx context.Context, args dto.ReadListUserParams) ([]model.User, error) {
	call := m.Called(ctx, args)

	var r0 []model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadCountUserQuery(ctx context.Context, args dto.ReadListUserParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadUserByIDQuery(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) ReadUserProfileQuery(ctx context.Context, id string) (model.User, error) {
	call := m.Called(ctx, id)

	var r0 model.User
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.User); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *UserRepositoryMock) IsUserEmailExistsQuery(ctx context.Context, email string) (bool, error) {
	call := m.Called(ctx, email)
	return call.Bool(0), call.Error(1)
}

func (m *UserRepositoryMock) IsUpdateUserEmailExistsQuery(ctx context.Context, email, id string) (bool, error) {
	call := m.Called(ctx, email, id)
	return call.Bool(0), call.Error(1)
}
