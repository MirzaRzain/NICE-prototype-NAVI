package menu

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	menuRepo "be-dashboard-nba/internal/infrastructure/repository/menu"
)

type useCase struct {
	db          db.DB
	newMenuRepo func(q db.Query) contractRepo.MenuRepository
}

func NewUseCase(db db.DB) contract.MenuUseCase {
	return &useCase{
		db:          db,
		newMenuRepo: menuRepo.NewRepository,
	}
}
