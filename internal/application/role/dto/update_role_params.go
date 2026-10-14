package dto

import "database/sql"

type UpdateRoleParams struct {
	RoleID      int
	Name        string
	Code        string
	Description sql.NullString
	UpdatedBy   string
}
