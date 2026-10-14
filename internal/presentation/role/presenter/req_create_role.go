package presenter

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Code        string  `json:"code" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty"`
}
