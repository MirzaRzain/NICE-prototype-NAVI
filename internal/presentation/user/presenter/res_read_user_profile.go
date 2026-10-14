package presenter

import (
	"be-dashboard-nba/internal/domain/model"
)

type ReadUserProfileUseCaseResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	Name     string  `json:"name"`
	ImgPath  string  `json:"img_path"`
	ImgName  string  `json:"img_name"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	RoleID   int     `json:"role_id"`
	Phone    *string `json:"phone"`
	Active   bool    `json:"active"`
}

func ToReadUserProfileUseCaseResponse(entity model.User) (response ReadUserProfileUseCaseResponse) {
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

	if entity.ImgName.Valid {
		imgName := entity.ImgName.String
		response.ImgName = imgName
	}

	if entity.ImgPath.Valid {
		imgPath := entity.ImgPath.String
		response.ImgPath = imgPath
	}

	return
}
