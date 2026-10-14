package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/user/dto"
	"be-dashboard-nba/internal/domain/model"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
)

func ToUpdateUserParams(req *userPresenter.UpdateUserRequest, userID, updatedBy string, existingUser model.User) dto.UpdateUserParams {
	params := dto.UpdateUserParams{
		ID:        userID,
		Name:      req.Name,
		FullName:  req.FullName,
		Email:     req.Email,
		UpdatedBy: updatedBy,
		RoleID:    req.RoleID,
		Active:    existingUser.Active,
		Phone:     existingUser.Phone,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}

	if req.Phone != nil {
		// Preserve prior behavior: empty string means NULL.
		if *req.Phone == "" {
			params.Phone = sql.NullString{Valid: false}
		} else {
			params.Phone = sql.NullString{String: *req.Phone, Valid: true}
		}
	}

	return params
}
