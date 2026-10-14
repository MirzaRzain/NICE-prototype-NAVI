package auth

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/jwt"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestLogoutUseCase(t *testing.T) {
	tests := []struct {
		name    string
		claims  *jwt.AccessTokenPayload
		ip      string
		setup   func(authRepo *mockrepo.AuthRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/logout",
			claims: &jwt.AccessTokenPayload{
				UserID:    "u1",
				SessionID: "s1",
			},
			ip: "127.0.0.1",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("DeleteSessionQuery", testifymock.Anything, "s1").Return(nil).Once()
				authRepo.On("CreateLoginRecord", testifymock.Anything, testifymock.Anything).Return(nil).Once()
			},
		},
		{
			name: "negative/delete-session-error",
			claims: &jwt.AccessTokenPayload{
				UserID:    "u1",
				SessionID: "s1",
			},
			ip: "127.0.0.1",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("DeleteSessionQuery", testifymock.Anything, "s1").Return(errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/login-record-error-ignored",
			claims: &jwt.AccessTokenPayload{
				UserID:    "u1",
				SessionID: "s1",
			},
			ip: "127.0.0.1",
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("DeleteSessionQuery", testifymock.Anything, "s1").Return(nil).Once()
				authRepo.On("CreateLoginRecord", testifymock.Anything, testifymock.Anything).Return(errors.New("log down")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRepo := mockrepo.NewAuthRepository()
			uc, _, cleanup := newTestUseCaseWithAuthRepo(t, authRepo)
			defer cleanup()

			tt.setup(authRepo)

			err := uc.LogoutUseCase(context.Background(), tt.claims, tt.ip)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
