package repository

import (
	"be-dashboard-nba/internal/application/menu/menu-permission/dto"
	"be-dashboard-nba/internal/domain/model"
	"context"
)

func (r *repository) ReadMenuPermissionListQuery(
	ctx context.Context,
	args dto.ReadMenuPermissionParams,
) (data []model.MenuPermission, err error) {

	const stmt = `
	SELECT
		id, menu_id, code, action_name
	FROM
		app_menu_permission
	WHERE
		(CASE WHEN $1::bool THEN(
			action_name ILIKE $2
		) ELSE TRUE END)
		 AND menu_id = $6
	ORDER BY
		(CASE WHEN $3 = 'action_name ASC' THEN action_name END)
		ASC,
		(CASE WHEN $3 = 'action_name DESC' THEN action_name END)
		DESC,
		(CASE WHEN $3 = 'code ASC' THEN code END)
		ASC,
		(CASE WHEN $3 = 'code DESC' THEN code END)
		DESC,
		(CASE WHEN $3 = 'id ASC' THEN id END)
		ASC,
		(CASE WHEN $3 = 'id DESC' THEN id END)
		DESC
	LIMIT $4
	OFFSET $5
	`
	rows, err := r.db.QueryContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
		args.MenuID,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var mp model.MenuPermission

		if err = rows.Scan(
			&mp.ID,
			&mp.MenuID,
			&mp.Code,
			&mp.ActionName,
		); err != nil {
			return
		}

		data = append(data, mp)
	}

	return
}

func (r *repository) ReadMenuPermissionByIdQuery(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error) {
	const stmt = `
		SELECT 
			id, menu_id, code, action_name
		FROM
			app_menu_permission
		WHERE 
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, MenuPermissionID).Scan(
		&data.ID,
		&data.MenuID,
		&data.Code,
		&data.ActionName,
	)

	return
}

func (r *repository) ReadMenuPermissionCount(ctx context.Context,
	args dto.ReadMenuPermissionParams) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_menu_permission
			WHERE
				(CASE WHEN $1::bool THEN(
				action_name ILIKE $2
				) ELSE TRUE END)
		 		AND menu_id = $3
		`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
		args.MenuID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return
}
