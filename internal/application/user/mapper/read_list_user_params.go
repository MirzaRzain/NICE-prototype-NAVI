package mapper

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/user/dto"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
)

func ToReadListUserParams(req *userPresenter.ReadUserRequest) dto.ReadListUserParams {
	req.Init()

	if req.Limit <= 0 {
		req.Limit = constant.DefaultLimit
	}

	if req.Page <= 0 {
		req.Page = constant.DefaultPage
	}

	// Preserve existing behavior: Offset comes from PaginationPayload.Init().
	return dto.ReadListUserParams{
		SetSearch: req.SetSearch,
		Search:    req.Search,
		Order:     req.Order,
		Offset:    req.Offset,
		Limit:     req.Limit,
	}
}
