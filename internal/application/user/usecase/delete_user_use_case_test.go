package user

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

func TestDeleteUserUseCase(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		deletedBy string
		setup     func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock)
		wantErr   error
	}{
		{
			name:      "positive/delete-user",
			userID:    "u1",
			deletedBy: "admin",
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "u1").Return(model.User{ID: "u1"}, nil).Once()
				userRepo.On("DeleteUserQuery", testifymock.Anything, "u1").Return(nil).Once()
			},
		},
		{
			name:      "negative/not-found",
			userID:    "missing",
			deletedBy: "admin",
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "missing").Return(model.User{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrUserIdNotFound,
		},
		{
			name:      "negative/delete-error",
			userID:    "u1",
			deletedBy: "admin",
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "u1").Return(model.User{ID: "u1"}, nil).Once()
				userRepo.On("DeleteUserQuery", testifymock.Anything, "u1").Return(errors.New("delete failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name:      "edge/forbidden-self-delete",
			userID:    "same",
			deletedBy: "same",
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "same").Return(model.User{ID: "same"}, nil).Once()
			},
			wantErr: constant.ErrForbiddenSelfDelete,
		},
		{
			name:      "edge/rollback-error-still-returns-original-error",
			userID:    "missing",
			deletedBy: "admin",
			setup: func(m sqlmock.Sqlmock, userRepo *mockrepo.UserRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(sql.ErrConnDone)

				userRepo.On("ReadUserByIDQuery", testifymock.Anything, "missing").Return(model.User{}, sql.ErrNoRows).Once()
			},
			wantErr: constant.ErrUserIdNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockrepo.NewUserRepository()
			uc, mock, cleanup := newTestUseCaseForDeleteUser(t, userRepo)
			defer cleanup()

			tt.setup(mock, userRepo)

			err := uc.DeleteUserUseCase(context.Background(), tt.userID, tt.deletedBy)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
