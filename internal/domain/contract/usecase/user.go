package usecase

import (
	"context"

	"be-dashboard-nba/internal/domain/model"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
)

type UserUseCase interface {
	ReadDetailUserUseCase(ctx context.Context, id string) (data model.User, err error)
	UpdateUserUseCase(ctx context.Context, request userPresenter.UpdateUserRequest, updatedBy string, userID string) (err error)
	ReadUsersUseCase(ctx context.Context, req userPresenter.ReadUserRequest) (data model.UserPaginationResponse, err error)
	CreateUserUseCase(ctx context.Context, req userPresenter.CreateUserRequest, userID string) (err error)
	DeleteUserUseCase(ctx context.Context, userID string, deletedBy string) (err error)
	ReadUserProfileUseCase(ctx context.Context, userID string) (data model.User, err error)
}
