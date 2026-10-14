package repository

import (
	db "be-dashboard-nba/internal/infrastructure/database"
)

type Query struct {
	db db.Query
}

func NewQuery(db db.Query) *Query {
	return &Query{
		db: db,
	}
}
