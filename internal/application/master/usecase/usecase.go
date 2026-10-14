package usecase

import (
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	masterRepo "be-dashboard-nba/internal/infrastructure/repository/master"
)

type useCase struct {
	db            db.DB
	newMasterRepo func(db.DB) contractRepo.MasterRepository
}

func NewUseCase(database db.DB) contract.MasterUseCase {
	return &useCase{db: database, newMasterRepo: masterRepo.NewRepository}
}
