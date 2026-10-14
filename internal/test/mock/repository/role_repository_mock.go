package mock

import (
	"context"

	"be-dashboard-nba/internal/application/role/dto"
	"be-dashboard-nba/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func NewRoleRepository() *RoleRepositoryMock { return &RoleRepositoryMock{} }

func (m *RoleRepositoryMock) ReadRolesQuery(ctx context.Context, args dto.ReadRolesParams) ([]model.Role, error) {
	call := m.Called(ctx, args)

	var r0 []model.Role
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Role); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRolesCount(ctx context.Context, args dto.ReadRolesParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRoleByIDQuery(ctx context.Context, roleID int) (model.Role, error) {
	call := m.Called(ctx, roleID)

	var r0 model.Role
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Role); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) CreateRoleQuery(ctx context.Context, params dto.CreateRoleParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *RoleRepositoryMock) UpdateRoleQuery(ctx context.Context, params dto.UpdateRoleParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *RoleRepositoryMock) DeleteRoleQuery(ctx context.Context, roleID int) error {
	call := m.Called(ctx, roleID)
	return call.Error(0)
}

func (m *RoleRepositoryMock) ReadRoleAccessQuery(ctx context.Context, args dto.ReadRoleAccessParams) ([]model.RoleAccessResponse, error) {
	call := m.Called(ctx, args)

	var r0 []model.RoleAccessResponse
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.RoleAccessResponse); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) ReadRoleAccessCount(ctx context.Context, args dto.ReadRoleAccessParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *RoleRepositoryMock) DeleteRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) error {
	call := m.Called(ctx, payload)
	return call.Error(0)
}

func (m *RoleRepositoryMock) CreateRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) error {
	call := m.Called(ctx, payload)
	return call.Error(0)
}
