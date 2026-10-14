package auth

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	authRepo "be-dashboard-nba/internal/infrastructure/repository/auth"
)

type useCase struct {
	db          db.DB
	newAuthRepo func(q db.Query) contractRepo.AuthRepository
}

func NewUseCase(db db.DB) contract.AuthUseCase {
	return &useCase{
		db:          db,
		newAuthRepo: authRepo.NewRepo,
	}
}
