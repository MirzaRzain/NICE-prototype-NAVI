package user

import (
	"testing"

	"be-dashboard-nba/internal/application/utils"
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	db "be-dashboard-nba/internal/infrastructure/database"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithUserAndRoleRepo(
	t *testing.T,
	userRepo *mockrepo.UserRepositoryMock,
	roleRepo *mockrepo.RoleRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newUserRepo: func(db.Query) contractRepo.UserRepository { return userRepo },
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForDeleteUser(t *testing.T, userRepo *mockrepo.UserRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newUserRepo: func(db.Query) contractRepo.UserRepository { return userRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForUpdateUser(
	t *testing.T,
	userRepo *mockrepo.UserRepositoryMock,
	roleRepo *mockrepo.RoleRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithUserAndRoleRepo(t, userRepo, roleRepo)
}
