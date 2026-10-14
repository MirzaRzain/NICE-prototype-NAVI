package mapper

import (
	"be-dashboard-nba/internal/application/role/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
)

func ToUpdateRoleMenuPermissionParam(item *rolePresenter.UpdateRoleAccessItem, roleID int) dto.UpdateRoleMenuPermission {
	return dto.UpdateRoleMenuPermission{
		MenuPermissionID: item.AccessID,
		RoleID:           roleID,
	}
}
