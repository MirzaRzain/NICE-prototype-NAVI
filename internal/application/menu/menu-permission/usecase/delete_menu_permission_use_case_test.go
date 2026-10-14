package menu_permission

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestDeleteMenuPermissionUseCase(t *testing.T) {
	tests := []struct {
		name             string
		menuPermissionID int
		setup            func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock)
		wantErr          error
	}{
		{
			name:             "positive/delete-menu-permission",
			menuPermissionID: 1,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
				menuPermissionRepo.On("DeleteMenuPermissionQuery", testifymock.Anything, 1).Return(nil).Once()
			},
		},
		{
			name:             "negative/not-found",
			menuPermissionID: 999,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 999).Return(model.MenuPermission{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuPermissionIdNotFound,
		},
		{
			name:             "negative/delete-error",
			menuPermissionID: 1,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
				menuPermissionRepo.On("DeleteMenuPermissionQuery", testifymock.Anything, 1).Return(errors.New("delete failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:             "edge/rollback-error",
			menuPermissionID: 1,
			setup: func(m sqlmock.Sqlmock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{}, errors.New("read failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, mock, cleanup := newTestUseCaseForDelete(t, menuPermissionRepo)
			defer cleanup()

			tt.setup(mock, menuPermissionRepo)

			err := uc.DeleteMenuPermissionUseCase(context.Background(), tt.menuPermissionID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
