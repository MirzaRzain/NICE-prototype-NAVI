package dto

import "database/sql"

type CreateMenuParams struct {
	ParentID    sql.NullInt32
	Name        string
	Description sql.NullString
	URL         sql.NullString
	Sort        int
	Group       string
	Icon        sql.NullString
	Active      bool
	Display     bool
	CreatedBy   string
}
