package user

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	roleRepo "be-dashboard-nba/internal/infrastructure/repository/role"
	userRepo "be-dashboard-nba/internal/infrastructure/repository/user"
)

type useCase struct {
	db db.DB

	newUserRepo func(q db.Query) contractRepo.UserRepository
	newRoleRepo func(q db.Query) contractRepo.RoleRepository
}

func NewUseCase(db db.DB) contract.UserUseCase {
	return &useCase{
		db:          db,
		newUserRepo: userRepo.NewRepository,
		newRoleRepo: roleRepo.NewRepository,
	}
}
