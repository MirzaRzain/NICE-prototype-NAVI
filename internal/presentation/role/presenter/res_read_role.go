package presenter

import "be-dashboard-nba/internal/domain/model"

type ReadRoleResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description"`
}

func ToReadRoleResponse(entity model.Role) (response ReadRoleResponse) {
	response.ID = entity.ID
	response.Name = entity.Name
	response.Code = entity.Code

	if entity.Description.Valid {
		response.Description = &entity.Description.String
	}

	return
}

func ToReadRoleListResponses(entities []model.Role) (responses []ReadRoleResponse) {
	responses = make([]ReadRoleResponse, len(entities))

	for i := range entities {
		responses[i] = ToReadRoleResponse(entities[i])
	}

	return
}
