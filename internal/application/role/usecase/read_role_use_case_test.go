package role

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/utils"
	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadRolesUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     rolePresenter.ReadRolesRequest
		setup   func(roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/returns-pagination",
			req:  rolePresenter.ReadRolesRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCount", testifymock.Anything, testifymock.Anything).Return(3, nil).Once()
				roleRepo.On("ReadRolesQuery", testifymock.Anything, testifymock.Anything).Return([]model.Role{{ID: 1}}, nil).Once()
			},
		},
		{
			name: "negative/count-error",
			req:  rolePresenter.ReadRolesRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCount", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/list-error",
			req:  rolePresenter.ReadRolesRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCount", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				roleRepo.On("ReadRolesQuery", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/default-limit-page",
			req:  rolePresenter.ReadRolesRequest{},
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRolesCount", testifymock.Anything, testifymock.Anything).Return(int64(0), nil).Once()
				roleRepo.On("ReadRolesQuery", testifymock.Anything, testifymock.Anything).Return([]model.Role{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, _, cleanup := newTestUseCaseForCreateRole(t, roleRepo)
			defer cleanup()

			tt.setup(roleRepo)

			_, err := uc.ReadRolesUseCase(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
