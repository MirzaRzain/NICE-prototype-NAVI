package role

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	menuPermissionRepo "be-dashboard-nba/internal/infrastructure/repository/menu/menu-permission"
	roleRepo "be-dashboard-nba/internal/infrastructure/repository/role"
)

type useCase struct {
	db db.DB

	newRoleRepo           func(q db.Query) contractRepo.RoleRepository
	newMenuPermissionRepo func(q db.Query) contractRepo.MenuPermissionRepository
}

func NewUseCase(db db.DB) contract.RoleUseCase {
	return &useCase{
		db:                    db,
		newRoleRepo:           roleRepo.NewRepository,
		newMenuPermissionRepo: menuPermissionRepo.NewRepository,
	}
}
