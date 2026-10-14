package repository

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth/dto"
	"be-dashboard-nba/internal/domain/model"
)

type AuthRepository interface {
	CreateSessionQuery(ctx context.Context, args dto.SessionParams) (err error)
	UpdateSessionQuery(ctx context.Context, args dto.SessionParams) (err error)
	DeleteSessionQuery(ctx context.Context, id string) (err error)
	ReadDetailSessionQuery(ctx context.Context, id string) (data model.Session, err error)
	ReadDetailUserByEmailQuery(ctx context.Context, email string) (data model.User, err error)
	CreateLoginAttemp(ctx context.Context, args dto.LoginAttempParams) (err error)
	CreateLoginRecord(ctx context.Context, args dto.LoginRecord) (err error)
	CheckPermissionQuery(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (bool, error)
	ReadDetailUserByIdQuery(ctx context.Context, id string) (data model.User, err error)
}
