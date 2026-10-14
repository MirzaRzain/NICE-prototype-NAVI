package auth

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/auth/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase(t *testing.T) {
	_ = os.Setenv("AUTH_ACCESS_TOKEN_SECRET_KEY", "secret")
	_ = os.Setenv("AUTH_ACCESS_TOKEN_EXPIRES", "24h")
	_ = os.Setenv("AUTH_REFRESH_TOKEN_SECRET_KEY", "secret")
	_ = os.Setenv("AUTH_REFRESH_TOKEN_EXPIRES", "48h")

	hashed, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		req   authPresenter.LoginRequest
		setup func(m sqlmock.Sqlmock, authRepo *mockrepo.AuthRepositoryMock)

		wantErr error
	}{
		{
			name: "positive/login",
			req: authPresenter.LoginRequest{
				Email:    "a@example.com",
				Password: "password123",
			},
			setup: func(m sqlmock.Sqlmock, authRepo *mockrepo.AuthRepositoryMock) {
				m.ExpectBegin()
				m.ExpectCommit()

				authRepo.On("ReadDetailUserByEmailQuery", testifymock.Anything, "a@example.com").
					Return(model.User{ID: "u1", Password: string(hashed)}, nil).Once()
				authRepo.On("CreateSessionQuery", testifymock.Anything, testifymock.Anything).Return(nil).Once()
				authRepo.On("CreateLoginRecord", testifymock.Anything, testifymock.Anything).Return(nil).Once()
			},
		},
		{
			name: "negative/wrong-email",
			req: authPresenter.LoginRequest{
				Email:    "missing@example.com",
				Password: "password123",
			},
			setup: func(m sqlmock.Sqlmock, authRepo *mockrepo.AuthRepositoryMock) {
				m.ExpectBegin()

				authRepo.On("ReadDetailUserByEmailQuery", testifymock.Anything, "missing@example.com").
					Return(model.User{}, sql.ErrNoRows).Once()
				authRepo.On("CreateLoginAttemp", testifymock.Anything, testifymock.Anything).Return(nil).Once()
			},
			wantErr: constant.ErrWrongEmailOrPassword,
		},
		{
			name: "negative/wrong-password",
			req: authPresenter.LoginRequest{
				Email:    "a@example.com",
				Password: "wrong",
			},
			setup: func(m sqlmock.Sqlmock, authRepo *mockrepo.AuthRepositoryMock) {
				m.ExpectBegin()

				authRepo.On("ReadDetailUserByEmailQuery", testifymock.Anything, "a@example.com").
					Return(model.User{ID: "u1", Password: string(hashed)}, nil).Once()
				authRepo.On("CreateLoginAttemp", testifymock.Anything, testifymock.Anything).Return(nil).Once()
			},
			wantErr: constant.ErrWrongEmailOrPassword,
		},
		{
			name: "negative/create-session-error",
			req: authPresenter.LoginRequest{
				Email:    "a@example.com",
				Password: "password123",
			},
			setup: func(m sqlmock.Sqlmock, authRepo *mockrepo.AuthRepositoryMock) {
				m.ExpectBegin()
				m.ExpectRollback()

				authRepo.On("ReadDetailUserByEmailQuery", testifymock.Anything, "a@example.com").
					Return(model.User{ID: "u1", Password: string(hashed)}, nil).Once()
				authRepo.On("CreateSessionQuery", testifymock.Anything, testifymock.Anything).Return(errors.New("insert failed")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRepo := mockrepo.NewAuthRepository()
			uc, mock, cleanup := newTestUseCaseWithAuthRepo(t, authRepo)
			defer cleanup()

			tt.setup(mock, authRepo)

			_, _, err := uc.LoginUseCase(context.Background(), tt.req, "ua", "127.0.0.1")
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
