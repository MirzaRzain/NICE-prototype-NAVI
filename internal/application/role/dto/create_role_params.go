package dto

import "database/sql"

type CreateRoleParams struct {
	Name        string
	Code        string
	Description sql.NullString
	CreatedBy   string
}
