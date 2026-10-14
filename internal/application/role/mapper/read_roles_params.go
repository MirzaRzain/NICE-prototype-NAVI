package mapper

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
)

func ToReadRolesParams(req *rolePresenter.ReadRolesRequest) dto.ReadRolesParams {
	req.Init()
	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	return dto.ReadRolesParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}
}
