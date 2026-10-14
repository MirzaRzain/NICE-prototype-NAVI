package role

import (
	"context"
	"errors"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/dto"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCreateRoleUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     rolePresenter.CreateRoleRequest
		setup   func(roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/create-role",
			req: rolePresenter.CreateRoleRequest{
				Name: "Admin",
				Code: "ADMIN",
			},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.
					On("CreateRoleQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateRoleParams)
						assert.Equal(t, "ADMIN", params.Code)
						assert.Equal(t, "Admin", params.Name)
						assert.Equal(t, "user-1", params.CreatedBy)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "negative/repo-error",
			req: rolePresenter.CreateRoleRequest{
				Name: "Admin",
				Code: "ADMIN",
			},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("CreateRoleQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/empty-fields-still-maps-created-by",
			req:  rolePresenter.CreateRoleRequest{},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.
					On("CreateRoleQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateRoleParams)
						assert.Equal(t, "user-1", params.CreatedBy)
					}).
					Return(nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, _, cleanup := newTestUseCaseForCreateRole(t, roleRepo)
			defer cleanup()

			tt.setup(roleRepo)

			err := uc.CreateRoleUseCase(context.Background(), tt.req, "user-1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
