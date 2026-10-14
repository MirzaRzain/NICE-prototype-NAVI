package user

import (
	"context"
	"testing"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/utils"
	"be-dashboard-nba/internal/domain/model"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
	mockrepo "be-dashboard-nba/internal/test/mock/repository"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

func TestReadUsersUseCase(t *testing.T) {
	tests := []struct {
		name    string
		req     userPresenter.ReadUserRequest
		setup   func(userRepo *mockrepo.UserRepositoryMock)
		wantErr error
	}{
		{
			name: "positive/returns-pagination",
			req:  userPresenter.ReadUserRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQuery", testifymock.Anything, testifymock.Anything).Return(3, nil).Once()
				userRepo.On("ReadUsersQuery", testifymock.Anything, testifymock.Anything).Return([]model.User{{ID: "u1"}}, nil).Once()
			},
		},
		{
			name: "negative/count-error",
			req:  userPresenter.ReadUserRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQuery", testifymock.Anything, testifymock.Anything).Return(0, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "negative/list-error",
			req:  userPresenter.ReadUserRequest{PaginationPayload: utils.PaginationPayload{Limit: 2, Page: 1}},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQuery", testifymock.Anything, testifymock.Anything).Return(1, nil).Once()
				userRepo.On("ReadUsersQuery", testifymock.Anything, testifymock.Anything).Return(nil, errors.New("db down")).Once()
			},
			wantErr: constant.ErrUnknownSource,
		},
		{
			name: "edge/default-limit-page",
			req:  userPresenter.ReadUserRequest{},
			setup: func(userRepo *mockrepo.UserRepositoryMock) {
				userRepo.On("ReadCountUserQuery", testifymock.Anything, testifymock.Anything).Return(0, nil).Once()
				userRepo.On("ReadUsersQuery", testifymock.Anything, testifymock.Anything).Return([]model.User{}, nil).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mockrepo.NewUserRepository()
			uc, _, cleanup := newTestUseCaseForDeleteUser(t, userRepo)
			defer cleanup()

			tt.setup(userRepo)

			_, err := uc.ReadUsersUseCase(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}
