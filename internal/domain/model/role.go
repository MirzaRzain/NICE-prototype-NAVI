package model

import (
	"database/sql"
	"time"

	shareddto "be-dashboard-nba/internal/application/shared/dto"
)

type Role struct {
	ID          int
	Code        string
	Name        string
	Description sql.NullString
	CreatedBy   sql.NullTime
	CreatedAt   time.Time
	UpdatedBy   sql.NullString
	UpdatedAt   sql.NullTime
}

type RoleAccessResponse struct {
	RoleID         int
	RoleName       string
	MenuID         int
	MenuName       string
	PermissionID   int
	PermissionName string
	PermissionCode string
	HasAccess      bool
}

type RolePaginationResponse struct {
	Data       []Role
	Pagination shareddto.Pagination
}

type RoleAccessPaginationResponse struct {
	Data       []RoleAccessResponse
	Pagination shareddto.Pagination
}
