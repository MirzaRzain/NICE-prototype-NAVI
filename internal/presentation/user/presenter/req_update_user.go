package presenter

type UpdateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	FullName string  `json:"full_name"  validate:"required"`
	Email    string  `json:"email"  validate:"required,email"`
	RoleID   int     `json:"role_id" validate:"required"`
	Phone    *string `json:"phone"`
	Active   *bool   `json:"active"  validate:"required"`
}
