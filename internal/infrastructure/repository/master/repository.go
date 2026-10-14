package repository

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	db "be-dashboard-nba/internal/infrastructure/database"
)

type repo struct{ db db.DB }

func NewRepository(query db.DB) contractRepo.MasterRepository { return &repo{db: query} }
