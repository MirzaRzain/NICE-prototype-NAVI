package dto

import "database/sql"

type UpdateUserParams struct {
	ID        string
	Name      string
	FullName  string
	Email     string
	Active    bool
	Phone     sql.NullString
	UpdatedBy string
	RoleID    int
}
