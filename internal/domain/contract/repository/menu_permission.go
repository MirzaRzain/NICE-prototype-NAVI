package repository

import (
	"context"

	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	"be-dashboard-nba/internal/domain/model"
)

type MenuPermissionRepository interface {
	ReadMenuPermissionListQuery(ctx context.Context, args dto.ReadMenuPermissionParams) (data []model.MenuPermission, err error)
	ReadMenuPermissionByIdQuery(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error)
	ReadMenuPermissionCount(ctx context.Context, args dto.ReadMenuPermissionParams) (count int, err error)
	UpdateMenuPermissionQuery(ctx context.Context, params dto.UpdateMenuPermissionParams) (err error)
	CreateMenuPermissionQuery(ctx context.Context, params dto.CreateMenuPermissionParams) (err error)
	DeleteMenuPermissionQuery(ctx context.Context, MenuPermissionID int) (err error)
}
