package role

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestUpdateRoleAccessUseCase(t *testing.T) {
	tests := []struct {
		name    string
		roleID  int
		req     rolePresenter.UpdateRoleAccessRequest
		setup   func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock, createCalls *int, deleteCalls *int)
		wantErr error
		assert  func(createCalls int, deleteCalls int)
	}{
		{
			name:   "positive/create-and-delete",
			roleID: 10,
			req: rolePresenter.UpdateRoleAccessRequest{
				AccessItem: []rolePresenter.UpdateRoleAccessItem{
					{AccessID: 1, HasAccess: ptrBool(true)},
					{AccessID: 2, HasAccess: ptrBool(false)},
				},
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock, createCalls *int, deleteCalls *int) {
				m.ExpectBegin()
				m.ExpectCommit()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Twice()
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{ID: 1}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 2).Return(model.MenuPermission{ID: 2}, nil).Once()

				roleRepo.
					On("CreateRoleAccess", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) { *createCalls++ }).
					Return(nil).
					Once()
				roleRepo.
					On("DeleteRoleAccess", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) { *deleteCalls++ }).
					Return(nil).
					Once()
			},
			assert: func(createCalls int, deleteCalls int) {
				assert.Equal(t, 1, createCalls)
				assert.Equal(t, 1, deleteCalls)
			},
		},
		{
			name:   "edge/empty-access-item-commits",
			roleID: 10,
			req:    rolePresenter.UpdateRoleAccessRequest{AccessItem: []rolePresenter.UpdateRoleAccessItem{}},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.RoleRepositoryMock, _ *mockrepo.MenuPermissionRepositoryMock, _ *int, _ *int) {
				m.ExpectBegin()
				m.ExpectCommit()
			},
			assert: func(createCalls int, deleteCalls int) {
				assert.Equal(t, 0, createCalls)
				assert.Equal(t, 0, deleteCalls)
			},
		},
		{
			name:   "negative/menu-permission-not-found",
			roleID: 10,
			req: rolePresenter.UpdateRoleAccessRequest{
				AccessItem: []rolePresenter.UpdateRoleAccessItem{
					{AccessID: 1, HasAccess: ptrBool(true)},
				},
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock, menuPermissionRepo *mockrepo.MenuPermissionRepositoryMock, _ *int, _ *int) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				menuPermissionRepo.On("ReadMenuPermissionByIdQuery", testifymock.Anything, 1).Return(model.MenuPermission{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrMenuPermissionIdNotFound,
		},
		{
			name:   "negative/role-not-found",
			roleID: 10,
			req: rolePresenter.UpdateRoleAccessRequest{
				AccessItem: []rolePresenter.UpdateRoleAccessItem{
					{AccessID: 1, HasAccess: ptrBool(true)},
				},
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock, _ *mockrepo.MenuPermissionRepositoryMock, _ *int, _ *int) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			menuPermissionRepo := mockrepo.NewMenuPermissionRepository()
			uc, mock, cleanup := newTestUseCaseForUpdateAccess(t, roleRepo, menuPermissionRepo)
			defer cleanup()

			createCalls := 0
			deleteCalls := 0
			tt.setup(mock, roleRepo, menuPermissionRepo, &createCalls, &deleteCalls)

			err := uc.UpdateRoleAccessUseCase(context.Background(), tt.roleID, tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
			if tt.assert != nil {
				tt.assert(createCalls, deleteCalls)
			}
		})
	}
}

func ptrBool(v bool) *bool { return &v }
