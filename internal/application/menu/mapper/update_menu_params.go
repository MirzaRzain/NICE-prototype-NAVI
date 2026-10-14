package mapper

import (
	"database/sql"
	"strings"

	"be-dashboard-nba/internal/application/menu/dto"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
)

func ToUpdateMenuParams(req *menuPresenter.UpdateMenuRequest, updatedBy string, menuID int) dto.UpdateMenuParams {
	params := dto.UpdateMenuParams{
		ID:      menuID,
		Name:    req.Name,
		Group:   req.Group,
		Active:  false,
		Display: false,
	}

	if req.Active != nil {
		params.Active = *req.Active
	}

	if req.Display != nil {
		params.Display = *req.Display
	}

	if updatedBy != "" {
		params.UpdatedBy = sql.NullString{String: updatedBy, Valid: true}
	}

	if req.ParentID == nil || *req.ParentID == 0 {
		params.ParentID = sql.NullInt32{Valid: false}
	} else {
		params.ParentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	}

	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	if req.URL != nil {
		urlToStore := *req.URL
		if !strings.HasPrefix(urlToStore, "/") {
			urlToStore = "/" + urlToStore
		}
		params.URL = sql.NullString{String: urlToStore, Valid: true}
	}

	if req.Icon != nil {
		params.Icon = sql.NullString{String: *req.Icon, Valid: true}
	}

	return params
}
