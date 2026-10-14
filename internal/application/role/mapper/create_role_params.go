package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/role/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
)

func ToCreateRoleParams(req *rolePresenter.CreateRoleRequest, createdBy string) dto.CreateRoleParams {
	params := dto.CreateRoleParams{
		Name:      req.Name,
		Code:      req.Code,
		CreatedBy: createdBy,
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	return params
}
