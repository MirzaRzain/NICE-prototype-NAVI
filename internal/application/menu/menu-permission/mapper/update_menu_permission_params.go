package mapper

import (
	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
)

func ToUpdateMenuPermissionParams(req *menuPermissionPresenter.UpdateMenuPermissionRequest, updatedBy string, menuPermissionID int) dto.UpdateMenuPermissionParams {
	return dto.UpdateMenuPermissionParams{
		Code:             req.Code,
		ActionName:       req.ActionName,
		UpdatedBy:        updatedBy,
		MenuPermissionID: menuPermissionID,
	}
}
