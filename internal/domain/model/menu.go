package model

import (
	"database/sql"
	"time"
)

type Menu struct {
	ID          int            `json:"id"`
	ParentID    sql.NullInt32  `json:"parent_id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	URL         sql.NullString `json:"url"`
	Sort        int            `json:"sort"`
	Group       string         `json:"group"`
	Icon        sql.NullString `json:"icon"`
	Active      bool           `json:"active"`
	Display     bool           `json:"display"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedBy   sql.NullString `json:"updated_by"`
	UpdatedAt   sql.NullTime   `json:"updated_at"`
}
