package repository

import (
	"context"

	"be-dashboard-nba/internal/application/role/dto"
)

func (r *repository) CreateRoleQuery(ctx context.Context, params dto.CreateRoleParams) (err error) {
	const stmt = `
		INSERT INTO app_role(
			code, name, description, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`

	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.Name,
		params.Description,
		params.CreatedBy,
	)

	return
}

func (r *repository) UpdateRoleQuery(ctx context.Context, params dto.UpdateRoleParams) (err error) {
	const stmt = `
		UPDATE app_role
		SET
			code = $1,
			name = $2,
			description = COALESCE($3, description),
			updated_by = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $5
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.Name,
		params.Description,
		params.UpdatedBy,
		params.RoleID,
	)
	return
}

func (r *repository) DeleteRoleQuery(ctx context.Context, roleID int) (err error) {
	const stmt = `
		DELETE FROM app_role
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, roleID)

	return
}

func (r *repository) CreateRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) (err error) {
	const stmt = `
		INSERT INTO app_role_access(
			role_id, menu_permission_id
		)
		VALUES ($1, $2)
		ON CONFLICT (role_id, menu_permission_id) DO NOTHING;
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}

func (r *repository) DeleteRoleAccess(ctx context.Context, payload dto.UpdateRoleMenuPermission) (err error) {
	const stmt = `
		DELETE FROM app_role_access
		WHERE 
			role_id = $1 AND menu_permission_id = $2		
	`
	_, err = r.db.ExecContext(ctx, stmt,
		payload.RoleID,
		payload.MenuPermissionID,
	)
	return
}
