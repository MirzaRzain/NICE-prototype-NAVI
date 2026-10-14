package model

import (
	"database/sql"
	"time"

	shareddto "be-dashboard-nba/internal/application/shared/dto"
)

type MenuPermission struct {
	ID         int            `json:"id"`
	MenuID     int            `json:"menu_id"`
	Code       string         `json:"code"`
	ActionName string         `json:"action_name"`
	CreatedBy  time.Time      `json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedBy  sql.NullString `json:"updated_by"`
	UpdatedAt  sql.NullTime   `json:"updated_at"`
}

type MenuPermissionPaginationResponse struct {
	Data       []MenuPermission
	Pagination shareddto.Pagination
}
