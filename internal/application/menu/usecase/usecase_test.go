package menu

import (
	"testing"

	"be-dashboard-nba/internal/application/utils"
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	db "be-dashboard-nba/internal/infrastructure/database"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithMenuRepo(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newMenuRepo: func(db.Query) contractRepo.MenuRepository { return menuRepo },
	}
	return uc, mock, cleanup
}

func newTestUseCaseForDeleteMenu(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}

func newTestUseCaseForReadMenu(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}

func newTestUseCaseForUpdateMenuOrder(t *testing.T, menuRepo *mockrepo.MenuRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()
	return newTestUseCaseWithMenuRepo(t, menuRepo)
}
