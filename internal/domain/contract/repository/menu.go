package repository

import (
	"context"

	"be-dashboard-nba/internal/application/menu/dto"
	"be-dashboard-nba/internal/domain/model"
)

type MenuRepository interface {
	ReadSidebarMenuQuery(ctx context.Context, userID string) (data []model.Menu, err error)
	ReadListMenuQuery(ctx context.Context, params dto.ReadListMenuParams) (data []model.Menu, err error)
	ReadCountMenuQuery(ctx context.Context, params dto.ReadListMenuParams) (count int64, err error)
	ReadParentMenuQuery(ctx context.Context) (data []model.Menu, err error)
	CreateMenuQuery(ctx context.Context, params dto.CreateMenuParams) (err error)
	UpdateMenuQuery(ctx context.Context, params dto.UpdateMenuParams) (err error)
	DeleteMenuQuery(ctx context.Context, menuID int) (err error)
	ReadMenuByIDQuery(ctx context.Context, menuID int) (data model.Menu, err error)
	ReadSortForGroup(ctx context.Context, group string) (int, error)
	ReadNextSortForParent(ctx context.Context, parentID int32) (int, error)
	UpdateMenuSortQuery(ctx context.Context, params dto.UpdateMenuSortParams) (err error)
	ReadNextSortForParentAndGroup(ctx context.Context, parentID int32, group string) (int, error)
	CountMenuChildren(ctx context.Context, menuID int) (int, error)
	UpdateChildrenGroup(ctx context.Context, parentID int, newGroup string) error
}
