package user

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/user/dto"
	"be-dashboard-nba/internal/domain/model"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCreateUserUseCase(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		req   userPresenter.CreateUserRequest
		setup func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock)

		wantErr error
	}{
		{
			name: "positive/create-user",
			env:  "4",
			req: userPresenter.CreateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				Password: "password123",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 1).Return(model.Role{ID: 1}, nil).Once()
				userRepo.
					On("CreateUserQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.CreateUserParams)
						assert.Equal(t, "user-1", params.CreatedBy)
						assert.Equal(t, 1, params.RoleID)
						assert.NotEmpty(t, params.Password)
					}).
					Return("new-user", nil).
					Once()
				userRepo.On("CreateUserRoleQuery", testifymock.Anything, 1, "new-user").Return(nil).Once()
			},
		},
		{
			name: "negative/role-not-found",
			env:  "4",
			req: userPresenter.CreateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				Password: "password123",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 1).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name: "negative/create-error",
			env:  "4",
			req: userPresenter.CreateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				Password: "password123",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 1).Return(model.Role{ID: 1}, nil).Once()
				userRepo.On("CreateUserQuery", testifymock.Anything, testifymock.Anything).Return("", errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/invalid-bcrypt-cost-env",
			env:  "x",
			req: userPresenter.CreateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				Password: "password123",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, _ *mockrepo.UserRepositoryMock, _ *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("AUTH_BCRYPT_COST", tt.env)

			userRepo := mockrepo.NewUserRepository()
			roleRepo := mockrepo.NewRoleRepository()
			uc, mock, cleanup := newTestUseCaseWithUserAndRoleRepo(t, userRepo, roleRepo)
			defer cleanup()

			tt.setup(mock, userRepo, roleRepo)

			err := uc.CreateUserUseCase(context.Background(), tt.req, "user-1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func ptrBool(v bool) *bool { return &v }
