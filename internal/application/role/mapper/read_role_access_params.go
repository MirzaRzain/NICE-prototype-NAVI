package mapper

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
)

func ToReadRoleAccessParams(req *rolePresenter.ReadRoleAccessesRequest, roleID int) dto.ReadRoleAccessParams {
	req.Init()
	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	return dto.ReadRoleAccessParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
		RoleID:    roleID,
	}
}
