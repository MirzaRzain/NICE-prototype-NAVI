package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/user/dto"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
)

func ToCreateUserParams(req *userPresenter.CreateUserRequest, createdBy, hashedPassword string) dto.CreateUserParams {
	params := dto.CreateUserParams{
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedBy: createdBy,
		RoleID:    req.RoleID,
	}

	if req.Active != nil {
		params.Active = *req.Active
	} else {
		params.Active = true
	}

	if req.Phone != nil {
		params.Phone = sql.NullString{String: *req.Phone, Valid: true}
	}

	return params
}
