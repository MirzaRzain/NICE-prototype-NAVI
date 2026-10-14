package auth

import (
	"testing"

	"be-dashboard-nba/internal/application/utils"
	contractRepo "be-dashboard-nba/internal/domain/contract/repository"
	db "be-dashboard-nba/internal/infrastructure/database"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
)

func newTestUseCaseWithAuthRepo(t *testing.T, authRepo *mockrepo.AuthRepositoryMock) (*useCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	mockDB, mock, cleanup := utils.NewMockDB(t)
	uc := &useCase{
		db:          mockDB,
		newAuthRepo: func(db.Query) contractRepo.AuthRepository { return authRepo },
	}
	return uc, mock, cleanup
}
