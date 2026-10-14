package presenter

type UpdateMenuOrderRequest struct {
	Group     string `json:"group" validate:"required"`
	ParentID  *int   `json:"parent_id" validate:"omitempty,min=1"`
	SortedIDs []int  `json:"sorted_ids" validate:"required,min=1,dive,min=1"`
}
