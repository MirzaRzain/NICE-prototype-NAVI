package dto

import "database/sql"

type CreateUserParams struct {
	Name      string
	FullName  string
	Email     string
	Password  string
	RoleID    int
	Active    bool
	Phone     sql.NullString
	CreatedBy string
}
