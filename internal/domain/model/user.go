package model

import (
	"database/sql"
	"time"

	shareddto "be-dashboard-nba/internal/application/shared/dto"
)

type User struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	FullName  string         `json:"full_name"`
	Role      string         `json:"role"`
	RoleID    int            `json:"role_id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Active    bool           `json:"active"`
	Phone     sql.NullString `json:"phone"`
	ImgPath   sql.NullString `json:"img_path"`
	ImgName   sql.NullString `json:"img_name"`
	CreatedBy string         `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedBy sql.NullString `json:"updated_by"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
}

type UserPaginationResponse struct {
	Data       []User
	Pagination shareddto.Pagination
}
