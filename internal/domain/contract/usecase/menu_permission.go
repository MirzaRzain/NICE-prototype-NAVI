package usecase

import (
	"context"

	"be-dashboard-nba/internal/domain/model"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
)

type MenuPermissionUseCase interface {
	ReadMenuPermissionUseCase(ctx context.Context, args menuPermissionPresenter.ReadMenuPermissionListRequest, MenuID int) (data model.MenuPermissionPaginationResponse, err error)
	CreateMenuPermissionUseCase(ctx context.Context, payload menuPermissionPresenter.CreateMenuPermissionRequest, userID string, menuID int) (err error)
	UpdateMenuPermissionUseCase(ctx context.Context, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) (err error)
	ReadMenuPermissionDetailUseCase(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error)
	DeleteMenuPermissionUseCase(ctx context.Context, menuPermissionID int) (err error)
}
