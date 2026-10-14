package role

import (
	"testing"

	"be-dashboard-nba/internal/application/utils"
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	db "be-dashboard-nba/internal/infrastructure/database"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseForUpdateRole(t *testing.T, roleRepo *mockrepo.RoleRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForCreateRole(t *testing.T, roleRepo *mockrepo.RoleRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	return newTestUseCaseForUpdateRole(t, roleRepo)
}

func newTestUseCaseForUpdateAccess(
	t *testing.T,
	roleRepo *mockrepo.RoleRepositoryMock,
	menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock,
) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newRoleRepo: func(db.Query) contractRepo.RoleRepository { return roleRepo },
		newMenuPermissionRepo: func(db.Query) contractRepo.MenuPermissionRepository {
			return menuPermissionRepo
		},
	}
	return uc, mock, cleanup
}
