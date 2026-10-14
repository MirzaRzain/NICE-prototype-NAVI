package role

import (
	"context"
	"database/sql"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/dto"
	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestUpdateRoleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		roleID  int
		req     rolePresenter.UpdateRoleRequest
		setup   func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/update-role",
			roleID: 10,
			req: rolePresenter.UpdateRoleRequest{
				Name: "Admin",
				Code: "ADMIN",
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.
					On("UpdateRoleQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.UpdateRoleParams)
						assert.Equal(t, 10, params.RoleID)
						assert.Equal(t, "user-1", params.UpdatedBy)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name:   "negative/not-found",
			roleID: 999,
			req: rolePresenter.UpdateRoleRequest{
				Name: "Admin",
				Code: "ADMIN",
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 999).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name:   "negative/update-error",
			roleID: 10,
			req: rolePresenter.UpdateRoleRequest{
				Name: "Admin",
				Code: "ADMIN",
			},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
				roleRepo.On("UpdateRoleQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("update failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:   "edge/rollback-error",
			roleID: 10,
			req:    rolePresenter.UpdateRoleRequest{Name: "Admin", Code: "ADMIN"},
			setup: func(m sqlmock.Sqlmock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{}, errors.New("read failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, mock, cleanup := newTestUseCaseForUpdateRole(t, roleRepo)
			defer cleanup()

			tt.setup(mock, roleRepo)

			err := uc.UpdateRoleUseCase(context.Background(), tt.req, "user-1", tt.roleID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
