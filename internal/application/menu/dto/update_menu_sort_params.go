package dto

import "database/sql"

type UpdateMenuSortParams struct {
	ID        int
	Sort      int
	UpdatedBy string
	ParentID  sql.NullInt32
	Group     string
}
