package repository

import (
	"context"

	"be-dashboard-nba/internal/application/user/dto"
	"be-dashboard-nba/internal/domain/model"
)

type UserRepository interface {
	CreateUserQuery(ctx context.Context, args dto.CreateUserParams) (userID string, err error)
	CreateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error)
	UpdateUserQuery(ctx context.Context, args dto.UpdateUserParams) (err error)
	UpdateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error)
	DeleteUserQuery(ctx context.Context, id string) (err error)
	ReadUsersQuery(ctx context.Context, args dto.ReadListUserParams) (data []model.User, err error)
	ReadCountUserQuery(ctx context.Context, args dto.ReadListUserParams) (count int, err error)
	ReadUserByIDQuery(ctx context.Context, id string) (data model.User, err error)
	ReadUserProfileQuery(ctx context.Context, id string) (data model.User, err error)
	IsUserEmailExistsQuery(ctx context.Context, email string) (exists bool, err error)
	IsUpdateUserEmailExistsQuery(ctx context.Context, email, id string) (exists bool, err error)
}
