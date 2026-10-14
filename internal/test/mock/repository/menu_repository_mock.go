package mock

import (
	"context"

	"be-dashboard-nba/internal/application/menu/dto"
	"be-dashboard-nba/internal/domain/model"

	"github.com/stretchr/testify/mock"
)

type MenuRepositoryMock struct {
	mock.Mock
}

func NewMenuRepository() *MenuRepositoryMock { return &MenuRepositoryMock{} }

func (m *MenuRepositoryMock) ReadSidebarMenuQuery(ctx context.Context, userID string) ([]model.Menu, error) {
	call := m.Called(ctx, userID)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadListMenuQuery(ctx context.Context, params dto.ReadListMenuParams) ([]model.Menu, error) {
	call := m.Called(ctx, params)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadCountMenuQuery(ctx context.Context, params dto.ReadListMenuParams) (int64, error) {
	call := m.Called(ctx, params)

	var r0 int64
	raw0 := call.Get(0)
	if casted, ok := raw0.(int64); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadParentMenuQuery(ctx context.Context) ([]model.Menu, error) {
	call := m.Called(ctx)

	var r0 []model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.([]model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) CreateMenuQuery(ctx context.Context, params dto.CreateMenuParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) UpdateMenuQuery(ctx context.Context, params dto.UpdateMenuParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) DeleteMenuQuery(ctx context.Context, menuID int) error {
	call := m.Called(ctx, menuID)
	return call.Error(0)
}

func (m *MenuRepositoryMock) ReadMenuByIDQuery(ctx context.Context, menuID int) (model.Menu, error) {
	call := m.Called(ctx, menuID)

	var r0 model.Menu
	raw0 := call.Get(0)
	if casted, ok := raw0.(model.Menu); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadSortForGroup(ctx context.Context, group string) (int, error) {
	call := m.Called(ctx, group)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) ReadNextSortForParent(ctx context.Context, parentID int32) (int, error) {
	call := m.Called(ctx, parentID)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) UpdateMenuSortQuery(ctx context.Context, params dto.UpdateMenuSortParams) error {
	call := m.Called(ctx, params)
	return call.Error(0)
}

func (m *MenuRepositoryMock) ReadNextSortForParentAndGroup(ctx context.Context, parentID int32, group string) (int, error) {
	call := m.Called(ctx, parentID, group)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) CountMenuChildren(ctx context.Context, menuID int) (int, error) {
	call := m.Called(ctx, menuID)

	var r0 int
	raw0 := call.Get(0)
	if casted, ok := raw0.(int); ok {
		r0 = casted
	}

	return r0, call.Error(1)
}

func (m *MenuRepositoryMock) UpdateChildrenGroup(ctx context.Context, parentID int, newGroup string) error {
	call := m.Called(ctx, parentID, newGroup)
	return call.Error(0)
}
