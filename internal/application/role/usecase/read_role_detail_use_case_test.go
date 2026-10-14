package role

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

func TestReadRoleDetailUseCase(t *testing.T) {
	tests := []struct {
		name    string
		roleID  int
		setup   func(roleRepo *mockrepo.RoleRepositoryMock)
		wantErr error
	}{
		{
			name:   "positive/returns-data",
			roleID: 10,
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{ID: 10}, nil).Once()
			},
		},
		{
			name:   "negative/not-found",
			roleID: 999,
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 999).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name:   "negative/repo-error",
			roleID: 10,
			setup: func(roleRepo *mockrepo.RoleRepositoryMock) {
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 10).Return(model.Role{}, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := mockrepo.NewRoleRepository()
			uc, _, cleanup := newTestUseCaseForCreateRole(t, roleRepo)
			defer cleanup()

			tt.setup(roleRepo)

			_, err := uc.ReadRoleDetailUseCase(context.Background(), tt.roleID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
