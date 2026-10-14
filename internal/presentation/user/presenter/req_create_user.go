package presenter

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required"`
	FullName string  `json:"full_name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	RoleID   int     `json:"role_id" validate:"required"`
	Phone    *string `json:"phone"`
	Active   *bool   `json:"active"  validate:"required"`
}
