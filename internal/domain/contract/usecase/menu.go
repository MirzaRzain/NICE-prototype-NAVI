package usecase

import (
	"context"

	"be-dashboard-nba/internal/domain/model"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
)

type MenuUseCase interface {
	CreateMenuUseCase(ctx context.Context, request menuPresenter.CreateMenuRequest, userID string) (err error)
	UpdateMenuUseCase(ctx context.Context, request menuPresenter.UpdateMenuRequest, userID string, MenuID int) (err error)
	ReadListMenuUseCase(ctx context.Context, request menuPresenter.ReadMenuListRequest) (data []model.Menu, err error)
	ReadSidebarMenuUseCase(ctx context.Context, userID string) (data []model.Menu, err error)
	DeleteMenuUseCase(ctx context.Context, menuID int) (err error)
	UpdateMenuOrderUseCase(ctx context.Context, request menuPresenter.UpdateMenuOrderRequest, userID string) (err error)
	ReadMenuDetailUseCase(ctx context.Context, menuID int) (data model.Menu, err error)
	ReadMenuParentUseCase(ctx context.Context) (data []model.Menu, err error)
}
