package repository

import (
	"context"

	"be-dashboard-nba/internal/application/role/dto"
	"be-dashboard-nba/internal/domain/model"
)

type RoleRepository interface {
	ReadRolesQuery(ctx context.Context, args dto.ReadRolesParams) (data []model.Role, err error)
	ReadRolesCount(ctx context.Context, args dto.ReadRolesParams) (count int, err error)
	ReadRoleByIDQuery(ctx context.Context, roleID int) (data model.Role, err error)
	CreateRoleQuery(ctx context.Context, params dto.CreateRoleParams) (err error)
	UpdateRoleQuery(ctx context.Context, params dto.UpdateRoleParams) (err error)
	DeleteRoleQuery(ctx context.Context, roleID int) (err error)
	ReadRoleAccessQuery(ctx context.Context, args dto.ReadRoleAccessParams) (data []model.RoleAccessResponse, err error)
	ReadRoleAccessCount(ctx context.Context, args dto.ReadRoleAccessParams) (count int, err error)
	DeleteRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) (err error)
	CreateRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) (err error)
}
