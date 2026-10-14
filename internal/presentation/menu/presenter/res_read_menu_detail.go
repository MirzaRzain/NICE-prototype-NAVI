package presenter

import "be-dashboard-nba/internal/domain/model"

type MenuDetail struct {
	ID          int     `json:"id"`
	ParentID    *int32  `json:"parent_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	Sort        int     `json:"sort"`
	Group       string  `json:"group"`
	Icon        *string `json:"icon"`
	Active      bool    `json:"active"`
	Display     bool    `json:"display"`
}

func ToReadMenuDetailResponse(entity model.Menu) (data MenuDetail) {
	data = MenuDetail{
		ID:      entity.ID,
		Name:    entity.Name,
		Group:   entity.Group,
		Sort:    entity.Sort,
		Active:  entity.Active,
		Display: entity.Display,
	}

	if entity.ParentID.Valid {
		data.ParentID = &entity.ParentID.Int32
	}
	if entity.Description.Valid {
		data.Description = &entity.Description.String
	}
	if entity.Icon.Valid {
		data.Icon = &entity.Icon.String
	}
	if entity.URL.Valid {
		data.URL = &entity.URL.String
	}

	return

}
