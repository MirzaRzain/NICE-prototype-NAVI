package usecase

import (
	"context"

	"be-dashboard-nba/internal/domain/model"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
)

type RoleUseCase interface {
	ReadRolesUseCase(ctx context.Context, args rolePresenter.ReadRolesRequest) (data model.RolePaginationResponse, err error)
	ReadRoleDetailUseCase(ctx context.Context, roleID int) (data model.Role, err error)
	CreateRoleUseCase(ctx context.Context, paylaod rolePresenter.CreateRoleRequest, userID string) (err error)
	UpdateRoleUseCase(ctx context.Context, payload rolePresenter.UpdateRoleRequest, userID string, roleID int) (err error)
	DeleteRoleUseCase(ctx context.Context, roleID int) (err error)
	ReadRoleAccessUseCase(ctx context.Context, args rolePresenter.ReadRoleAccessesRequest, roleID int) (data model.RoleAccessPaginationResponse, err error)
	UpdateRoleAccessUseCase(ctx context.Context, roleID int, request rolePresenter.UpdateRoleAccessRequest) (err error)
}
