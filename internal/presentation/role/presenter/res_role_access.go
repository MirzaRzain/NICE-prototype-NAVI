package presenter

import "be-dashboard-nba/internal/domain/model"

type Access struct {
	AccessID   int    `json:"access_id"`
	AccessName string `json:"access_name"`
	HasAccess  bool   `json:"has_access"`
}

type RoleAccessResponse struct {
	MenuID   int      `json:"menu_id"`
	MenuName string   `json:"menu_name"`
	Accesses []Access `json:"accesses"`
}

func ToReadRoleAccessResponse(rows []model.RoleAccessResponse) []RoleAccessResponse {
	menuMap := make(map[int]*RoleAccessResponse)
	orderedMenuIDs := make([]int, 0)

	for _, r := range rows {
		if _, exists := menuMap[r.MenuID]; !exists {
			menuMap[r.MenuID] = &RoleAccessResponse{
				MenuID:   r.MenuID,
				MenuName: r.MenuName,
				Accesses: []Access{},
			}
			orderedMenuIDs = append(orderedMenuIDs, r.MenuID)
		}

		menuMap[r.MenuID].Accesses = append(menuMap[r.MenuID].Accesses, Access{
			AccessID:   r.PermissionID,
			AccessName: r.PermissionName,
			HasAccess:  r.HasAccess,
		})
	}

	result := make([]RoleAccessResponse, 0, len(menuMap))
	for _, id := range orderedMenuIDs {
		result = append(result, *menuMap[id])
	}

	return result
}
