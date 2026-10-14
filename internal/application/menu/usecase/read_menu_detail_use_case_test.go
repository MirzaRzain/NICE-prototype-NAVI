package menu

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadMenuDetailUseCase(t *testing.T) {
	tests := []struct {
		name    string
		menuID  int
		setup   func(menuRepo *mockrepo.MenuRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-menu",
			menuID: 1,
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			menuID: 999,
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 999).Return(model.Menu{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "negative/repo-error",
			menuID: 1,
			setup: func(menuRepo *mockrepo.MenuRepositoryMock) {
				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{}, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			uc, _, cleanup := newTestUseCaseForReadMenu(t, menuRepo)
			defer cleanup()

			tt.setup(menuRepo)

			_, err := uc.ReadMenuDetailUseCase(context.Background(), tt.menuID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
