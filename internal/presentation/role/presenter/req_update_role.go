package presenter

type UpdateRoleRequest struct {
	Name        string  `json:"name" validate:"required"`
	Code        string  `json:"code" validate:"required"`
	Description *string `json:"description" validate:"omitempty"`
}
