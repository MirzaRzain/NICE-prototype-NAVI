package mapper

import (
	"database/sql"

	"be-dashboard-nba/internal/application/menu/dto"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
)

func ToUpdateMenuSortParamsList(req *menuPresenter.UpdateMenuOrderRequest, updatedBy string) []dto.UpdateMenuSortParams {
	paramsList := make([]dto.UpdateMenuSortParams, len(req.SortedIDs))

	var parentID sql.NullInt32
	if req.ParentID != nil {
		parentID = sql.NullInt32{Int32: int32(*req.ParentID), Valid: true}
	} else {
		parentID = sql.NullInt32{Valid: false}
	}

	for index, menuID := range req.SortedIDs {
		paramsList[index] = dto.UpdateMenuSortParams{
			ID:        menuID,
			Sort:      index,
			UpdatedBy: updatedBy,
			Group:     req.Group,
			ParentID:  parentID,
		}
	}

	return paramsList
}
