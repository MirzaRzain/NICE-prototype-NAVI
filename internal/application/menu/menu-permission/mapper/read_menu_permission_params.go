package mapper

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
)

func ToReadMenuPermissionParams(req *menuPermissionPresenter.ReadMenuPermissionListRequest, menuID int) dto.ReadMenuPermissionParams {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	return dto.ReadMenuPermissionParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		MenuID:    menuID,
	}
}
