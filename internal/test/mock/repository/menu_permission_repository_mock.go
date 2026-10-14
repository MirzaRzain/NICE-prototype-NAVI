package mock

import (
	"context"

	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	"be-dashboard-nba/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MenuPermissionRepositoryMock struct {
	mock.Mock
}

func NewMenuPermissionRepository() *MenuPermissionRepositoryMock {
	return &MenuPermissionRepositoryMock{}
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionListQuery(ctx context.Context, args dto.ReadMenuPermissionParams) ([]model.MenuPermission, error) {
	call := m.Called(ctx, args)

	var r0 []model.MenuPermission
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.MenuPermission); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionByIdQuery(ctx context.Context, id int) (model.MenuPermission, error) {
	call := m.Called(ctx, id)

	var r0 model.MenuPermission
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.MenuPermission); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) ReadMenuPermissionCount(ctx context.Context, args dto.ReadMenuPermissionParams) (int, error) {
	call := m.Called(ctx, args)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuPermissionRepositoryMock) UpdateMenuPermissionQuery(ctx context.Context, params dto.UpdateMenuPermissionParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuPermissionRepositoryMock) CreateMenuPermissionQuery(ctx context.Context, params dto.CreateMenuPermissionParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuPermissionRepositoryMock) DeleteMenuPermissionQuery(ctx context.Context, id int) error {
	call := m.Called(ctx, id)
	return call.Error(0)
}
