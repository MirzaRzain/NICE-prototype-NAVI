package presenter

import (
	"be-dashboard-nba/internal/domain/model"
	"sort"
)

type MenuListItem struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Icon     *string        `json:"icon"`
	Url      *string        `json:"url"`
	Children []MenuListItem `json:"children"`
	Sort     int            `json:"sort"`
}

type ReadMenuListResponse struct {
	GroupName   string         `json:"group_name"`
	GroupChilds []MenuListItem `json:"group_childs"`
}

func ToReadMenuListResponse(menuEntities []model.Menu) []ReadMenuListResponse {
	if len(menuEntities) == 0 {
		return make([]ReadMenuListResponse, 0)
	}

	menuMap := make(map[int]model.Menu)
	for _, menu := range menuEntities {
		menuMap[menu.ID] = menu
	}

	rootMenus := make(map[string][]MenuListItem)

	for _, menu := range menuEntities {
		if !menu.ParentID.Valid {
			group := menu.Group
			rootMenus[group] = append(rootMenus[group], buildMenuTree(menu, menuMap))
		}
	}

	var response []ReadMenuListResponse
	for groupName, menus := range rootMenus {
		sortMenuItems(menus)
		response = append(response, ReadMenuListResponse{
			GroupName:   groupName,
			GroupChilds: menus,
		})
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].GroupName < response[j].GroupName
	})

	return response
}

func buildMenuTree(menu model.Menu, menuMap map[int]model.Menu) MenuListItem {
	item := MenuListItem{
		ID:       menu.ID,
		Name:     menu.Name,
		Sort:     menu.Sort,
		Children: []MenuListItem{},
	}

	if menu.Icon.Valid {
		item.Icon = &menu.Icon.String
	}
	if menu.URL.Valid {
		item.Url = &menu.URL.String
	}

	for _, potentialChild := range menuMap {
		if potentialChild.ParentID.Valid && int(potentialChild.ParentID.Int32) == menu.ID {
			item.Children = append(item.Children, buildMenuTree(potentialChild, menuMap))
		}
	}

	sortMenuItems(item.Children)

	return item
}
func sortMenuItems(items []MenuListItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sort < items[j].Sort
	})
	for i := range items {
		if len(items[i].Children) > 0 {
			sortMenuItems(items[i].Children)
		}
	}
}
