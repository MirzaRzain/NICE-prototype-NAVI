package menu_permission

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	menuRepo "be-dashboard-nba/internal/infrastructure/repository/menu"
	menuPermissionRepo "be-dashboard-nba/internal/infrastructure/repository/menu/menu-permission"
)

type useCase struct {
	db                    db.DB
	newMenuRepo           func(q db.Query) contractRepo.MenuRepository
	newMenuPermissionRepo func(q db.Query) contractRepo.MenuPermissionRepository
}

func NewUseCase(db db.DB) contract.MenuPermissionUseCase {
	return &useCase{
		db:                    db,
		newMenuRepo:           menuRepo.NewRepository,
		newMenuPermissionRepo: menuPermissionRepo.NewRepository,
	}
}
