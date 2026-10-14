package menu_permission

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	"be-dashboard-nba/internal/domain/model"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCreateMenuPermissionUseCase(t *testing.T) {
	tests := []struct {
		name    string
		menuID  int
		req     menuPermissionPresenter.CreateMenuPermissionRequest
		setup   func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/create-menu-permission",
			menuID: 1,
			req: menuPermissionPresenter.CreateMenuPermissionRequest{
				Code:       "R",
				ActionName: "read",
			},
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.
					On("CreateMenuPermissionQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateMenuPermissionParams)
						assert.Equal(t, "R", params.Code)
						assert.Equal(t, "read", params.ActionName)
						assert.Equal(t, 1, params.MenuID)
						assert.Equal(t, "user-1", params.CreatedBy)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name:   "negative/menu-not-found",
			menuID: 999,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, _ *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 999).Return(model.Menu{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuIdNotFound,
		},
		{
			name:   "negative/create-error",
			menuID: 1,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{ID: 1}, nil).Once()
				menuPermissionRepo.On("CreateMenuPermissionQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "edge/rollback-error",
			menuID: 1,
			setup: func(m sqlmock.Sqlmock, menuRepo *mockrepo.MenuRepositoryMock, _ *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuRepo.On("ReadMenuByIDQuery", testifymock.Anything, 1).Return(model.Menu{}, errors.New("read failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuRepo := mockrepo.NewMenuRepository()
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, mock, cleanup := newTestUseCaseWithRepos(t, menuRepo, menuPermissionRepo)
			defer cleanup()

			tt.setup(mock, menuRepo, menuPermissionRepo)

			err := uc.CreateMenuPermissionUseCase(context.Background(), tt.req, "user-1", tt.menuID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
