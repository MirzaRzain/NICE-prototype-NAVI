package presenter

import (
	"be-dashboard-nba/internal/domain/model"
)

type ReadUserDetailResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
}

func ToReadUserDetailResponse(entity model.User) (response ReadUserDetailResponse) {
	response.ID = entity.ID
	response.FullName = entity.FullName
	response.Name = entity.Name
	response.Email = entity.Email
	response.Role = entity.Role
	response.RoleID = entity.RoleID
	response.Active = entity.Active

	if entity.Phone.Valid {
		phone := entity.Phone.String
		response.Phone = &phone
	}

	return
}
