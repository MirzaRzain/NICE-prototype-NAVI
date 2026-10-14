package presenter

type UpdateMenuRequest struct {
	ParentID    *int    `json:"parent_id" validate:"omitempty,min=1"`
	Name        string  `json:"name" validate:"required,min=1,max=50"`
	Description *string `json:"description" validate:"omitempty,max=100"`
	URL         *string `json:"url" validate:"omitempty,max=100,uri"`
	Group       string  `json:"group" validate:"required,min=1,max=50"`
	Icon        *string `json:"icon" validate:"omitempty,max=50"`
	Active      *bool   `json:"active" validate:"required,boolean"`
	Display     *bool   `json:"display" validate:"required,boolean"`
}
