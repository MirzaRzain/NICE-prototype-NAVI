package auth

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestCheckPermissionUseCase(t *testing.T) {
	tests := []struct {
		name    string
		menuURL constant.MenuKey
		userID  string
		codes   []string
		setup   func(authRepo *mockrepo.AuthRepositoryMock)
		want    bool
		wantErr bool
	}{
		{
			name:    "positive/has-access",
			menuURL: constant.MenuSettingsUser,
			userID:  "u1",
			codes:   []string{"R"},
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("CheckPermissionQuery", testifymock.Anything, constant.MenuSettingsUser, "u1", []string{"R"}).Return(true, nil).Once()
			},
			want: true,
		},
		{
			name:    "negative/repo-error",
			menuURL: constant.MenuSettingsUser,
			userID:  "u1",
			codes:   []string{"R"},
			setup: func(authRepo *mockrepo.AuthRepositoryMock) {
				authRepo.On("CheckPermissionQuery", testifymock.Anything, constant.MenuSettingsUser, "u1", []string{"R"}).Return(false, errors.New("db down")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRepo := mockrepo.NewAuthRepository()
			uc, _, cleanup := newTestUseCaseWithAuthRepo(t, authRepo)
			defer cleanup()

			tt.setup(authRepo)

			hasAccess, err := uc.CheckPermissionUseCase(context.Background(), tt.menuURL, tt.userID, tt.codes)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, hasAccess)
		})
	}
}
