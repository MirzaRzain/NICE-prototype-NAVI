package repository

import (
	db "be-dashboard-nba/internal/infrastructure/database"

	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
)

type repository struct {
	db db.Query
}

func NewRepository(db db.Query) contractRepo.MenuPermissionRepository {
	return &repository{
		db: db,
	}
}
