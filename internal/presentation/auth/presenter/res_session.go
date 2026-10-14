package presenter

import "be-dashboard-nba/internal/domain/model"

type UserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Active   bool    `json:"active"`
	Phone    *string `json:"phone"`
	ImgPath  *string `json:"img_path"`
	ImgName  *string `json:"img_name"`
}

type SessionResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user,omitempty"`
}

func ToReadAuthMeResponse(entity model.User) (response UserResponse) {
	response.ID = entity.ID
	response.Name = entity.Name
	response.FullName = entity.FullName
	response.Email = entity.Email
	response.Active = entity.Active

	if entity.Phone.Valid {
		response.Phone = &entity.Phone.String
	}
	if entity.ImgPath.Valid {
		response.ImgPath = &entity.ImgPath.String
	}
	if entity.ImgName.Valid {
		response.ImgName = &entity.ImgName.String
	}

	return
}

func ToSessionResponse(entity model.Session, user model.User) (response SessionResponse) {
	response.AccessToken = entity.AccessToken
	response.RefreshToken = entity.RefreshToken
	if user.ID == "" {
		return
	}

	response.User = &UserResponse{}
	response.User.ID = user.ID
	response.User.Name = user.Name
	response.User.FullName = user.FullName
	response.User.Email = user.Email
	response.User.Active = user.Active

	if user.Phone.Valid {
		response.User.Phone = &user.Phone.String
	}
	if user.ImgPath.Valid {
		response.User.ImgPath = &user.ImgPath.String
	}
	if user.ImgName.Valid {
		response.User.ImgName = &user.ImgName.String
	}

	return
}
