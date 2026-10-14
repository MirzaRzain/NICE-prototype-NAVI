package usecase

import (
	"context"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/jwt"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/auth/presenter"
)

type AuthUseCase interface {
	LoginUseCase(ctx context.Context, request authPresenter.LoginRequest, userAgent, iPAddress string) (data model.Session, user model.User, err error)
	LogoutUseCase(ctx context.Context, claims *jwt.AccessTokenPayload, iPAddress string) (err error)
	AuthMeUseCase(ctx context.Context, id string) (data model.User, err error)
	CheckPermissionUseCase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error)
}
