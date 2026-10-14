package mapper

import (
	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
)

func ToCreateMenuPermissionParams(req *menuPermissionPresenter.CreateMenuPermissionRequest, createdBy string, menuID int) dto.CreateMenuPermissionParams {
	return dto.CreateMenuPermissionParams{
		Code:       req.Code,
		ActionName: req.ActionName,
		MenuID:     menuID,
		CreatedBy:  createdBy,
	}
}
