package menu_permission

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

func TestReadMenuPermissionDetailUseCase(t *testing.T) {
	tests := []struct {
		name             string
		menuPermissionID int
		setup            func(menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock)
		wantErr          error
	}{
		{
			name:             "positive/returns-data",
			menuPermissionID: 1,
			setup: func(menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
			},
		},
		{
			name:             "negative/not-found",
			menuPermissionID: 999,
			setup: func(menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 999).Return(model.MenuPermission{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuPermissionIdNotFound,
		},
		{
			name:             "negative/repo-error",
			menuPermissionID: 1,
			setup: func(menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{}, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, _, cleanup := newTestUseCaseForUpdate(t, menuPermissionRepo)
			defer cleanup()

			tt.setup(menuPermissionRepo)

			_, err := uc.ReadMenuPermissionDetailUseCase(context.Background(), tt.menuPermissionID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
