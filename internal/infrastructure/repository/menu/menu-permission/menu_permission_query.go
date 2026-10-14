package repository

import (
	"context"

	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
)

func (r *repository) CreateMenuPermissionQuery(ctx context.Context, params dto.CreateMenuPermissionParams) (err error) {
	const stmt = `
		INSERT INTO app_menu_permission (
			code, action_name, menu_id, created_by, created_at
		)
		VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`

	_, err = r.db.ExecContext(ctx, stmt, params.Code, params.ActionName, params.MenuID, params.CreatedBy)

	return
}

func (r *repository) UpdateMenuPermissionQuery(ctx context.Context, params dto.UpdateMenuPermissionParams) (err error) {
	const stmt = `
		UPDATE app_menu_permission
		SET
			code = $1,
			action_name = $2,
			updated_by = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $4
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Code,
		params.ActionName,
		params.UpdatedBy,
		params.MenuPermissionID,
	)

	return
}

func (r *repository) DeleteMenuPermissionQuery(ctx context.Context, MenuPermissionID int) (err error) {
	const stmt = `
		DELETE FROM app_menu_permission
		WHERE
			id = $1
	`
	_, err = r.db.ExecContext(ctx, stmt, MenuPermissionID)

	return
}
