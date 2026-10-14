package repository

import (
	db "be-dashboard-nba/internal/infrastructure/database"

	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
)

type repository struct {
	DB db.Query
}

func NewRepo(db db.Query) contractRepo.AuthRepository {
	return &repository{
		DB: db,
	}
}
