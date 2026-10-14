package presenter

import "be-dashboard-nba/internal/domain/model"

type ReadMenuPermissionResponse struct {
	ID         int    `json:"id"`
	MenuID     int    `json:"menu_id"`
	Code       string `json:"code"`
	ActionName string `json:"action_name"`
}

func ToReadMenuPermissionListResponse(entity model.MenuPermission) (response ReadMenuPermissionResponse) {
	response.ID = entity.ID
	response.MenuID = entity.MenuID
	response.Code = entity.Code
	response.ActionName = entity.ActionName

	return
}

func ToReadMenuPermissionListResponses(entities []model.MenuPermission) (responses []ReadMenuPermissionResponse) {
	responses = make([]ReadMenuPermissionResponse, len(entities))

	for i := range entities {
		responses[i] = ToReadMenuPermissionListResponse(entities[i])
	}

	return
}
