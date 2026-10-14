package user

import (
	"context"
	"database/sql"
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

func TestUpdateUserUseCase(t *testing.T) {
	tests := []struct {
		name  string
		req   userPresenter.UpdateUserRequest
		setup func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock)

		wantErr error
	}{
		{
			name: "positive/role-changed",
			req: userPresenter.UpdateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				RoleID:   2,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "user-99").Return(model.User{ID: "user-99", Role: "OLD", Active: true}, nil).Once()
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 2).Return(model.Role{ID: 2, Name: "NEW"}, nil).Once()
				userRepo.On("UpdateUserRoleQuery", testifymock.Anything, 2, "user-99").Return(nil).Once()
				userRepo.
					On("UpdateUserQuery", testifymock.Anything, testifymock.Anything).
					Run(func(args testifymock.Arguments) {
						params := args.Get(1).(dto.UpdateUserParams)
						assert.Equal(t, "user-99", params.ID)
						assert.Equal(t, 2, params.RoleID)
					}).
					Return(nil).
					Once()
			},
		},
		{
			name: "edge/role-unchanged-skips-update-user-role",
			req: userPresenter.UpdateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "user-99").Return(model.User{ID: "user-99", Role: "OLD", Active: true}, nil).Once()
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 1).Return(model.Role{ID: 1, Name: "OLD"}, nil).Once()
				userRepo.On("UpdateUserQuery", testifymock.Anything, testifymock.Anything).Return(nil).Once()
			},
		},
		{
			name: "negative/user-not-found",
			req: userPresenter.UpdateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, _ *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "user-99").Return(model.User{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrUserIdNotFound,
		},
		{
			name: "negative/role-not-found",
			req: userPresenter.UpdateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				RoleID:   2,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "user-99").Return(model.User{ID: "user-99", Role: "OLD", Active: true}, nil).Once()
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 2).Return(model.Role{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrRoleIdNotFound,
		},
		{
			name: "negative/update-error",
			req: userPresenter.UpdateUserRequest{
				Name:     "A",
				FullName: "A A",
				Email:    "a@example.com",
				RoleID:   1,
				Active:   ptrBool(true),
			},
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock, roleRepo *mockrepo.RoleRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "user-99").Return(model.User{ID: "user-99", Role: "OLD", Active: true}, nil).Once()
				roleRepo.On("ReadRoleByIDQuery", testifymock.Anything, 1).Return(model.Role{ID: 1, Name: "OLD"}, nil).Once()
				userRepo.On("UpdateUserQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("update failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockrepo.NewUserRepository()
			roleRepo := mockrepo.NewRoleRepository()
			uc, mock, cleanup := newTestUseCaseForUpdateUser(t, userRepo, roleRepo)
			defer cleanup()

			tt.setup(mock, userRepo, roleRepo)

			err := uc.UpdateUserUseCase(context.Background(), tt.req, "admin", "user-99")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)

			if tt.name == "edge/role-unchanged-skips-update-user-role" {
				userRepo.AssertNotCalled(t, "UpdateUserRoleQuery", testifymock.Anything, testifymock.Anything, testifymock.Anything)
			}
		})
	}
}
