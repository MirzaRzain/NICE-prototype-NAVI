package repository

import (
	"be-dashboard-nba/internal/application/role/dto"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"fmt"
)

func (r *repository) ReadRolesQuery(
	ctx context.Context,
	args dto.ReadRolesParams,
) (data []model.Role, err error) {
	const stmt = `
		SELECT 
			id, name, code ,description
		FROM
			app_role
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'name ASC' THEN name END)
			ASC,
			(CASE WHEN $3 = 'name DESC' THEN name END)
			DESC,
			(CASE WHEN $3 = 'description ASC' THEN description END)
			ASC,
			(CASE WHEN $3 = 'description DESC' THEN description END)
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
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var r model.Role

		if err = rows.Scan(
			&r.ID,
			&r.Name,
			&r.Code,
			&r.Description,
		); err != nil {
			return
		}

		data = append(data, r)
	}

	return
}

func (r *repository) ReadRolesCount(ctx context.Context,
	args dto.ReadRolesParams) (count int, err error) {
	const stmt = `
			SELECT COUNT(*)
			FROM
				app_role
			WHERE
				(CASE WHEN $1::bool THEN(
				name ILIKE $2
				) ELSE TRUE END)
		`
	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)
	if err != nil {
		fmt.Printf("Error in ReadRolesCount: %v\n", err)
	}
	return
}

func (r *repository) ReadRoleByIDQuery(ctx context.Context, roleID int) (data model.Role, err error) {
	const stmt = `
		SELECT 
			id, name, code, description
		FROM
			app_role
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, roleID).Scan(
		&data.ID,
		&data.Name,
		&data.Code,
		&data.Description,
	)

	return
}

type ReadRoleAccessParams = dto.ReadRoleAccessParams

func (r *repository) ReadRoleAccessQuery(
	ctx context.Context,
	args ReadRoleAccessParams,
) (data []model.RoleAccessResponse, err error) {
	const stmt = `
		WITH filtered_menu AS (
			SELECT 
				id, name
			FROM app_menu
			WHERE
				(CASE WHEN $2::bool THEN (
					name ILIKE $3
				) ELSE TRUE END)
			ORDER BY
				(CASE WHEN $4 = 'name ASC' THEN name END) ASC,
				(CASE WHEN $4 = 'name DESC' THEN name END) DESC,
				(CASE WHEN $4 = 'id ASC' THEN id END) ASC,
				(CASE WHEN $4 = 'id DESC' THEN id END) DESC,
				name ASC
			LIMIT $5
			OFFSET $6
		)
		SELECT 
			ar.id AS role_id,
			ar.name AS role_name,
			fm.id AS menu_id,
			fm.name AS menu_name,
			amp.id AS permission_id,
			amp.action_name AS permission_name,
			amp.code AS permission_code,
			CASE 
				WHEN ara.role_id IS NOT NULL THEN TRUE
				ELSE FALSE
			END AS has_access
		FROM filtered_menu fm
		JOIN app_menu_permission amp 
			ON fm.id = amp.menu_id
		JOIN app_role ar
			ON ar.id = $1
		LEFT JOIN app_role_access ara 
			ON amp.id = ara.menu_permission_id
			AND ara.role_id = $1
		ORDER BY fm.name ASC, amp.action_name ASC
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		args.RoleID,
		args.SetSearch,
		args.Search,
		args.Order,
		args.Limit,
		args.Offset,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var res model.RoleAccessResponse
		if err = rows.Scan(
			&res.RoleID,
			&res.RoleName,
			&res.MenuID,
			&res.MenuName,
			&res.PermissionID,
			&res.PermissionName,
			&res.PermissionCode,
			&res.HasAccess,
		); err != nil {
			return
		}
		data = append(data, res)
	}

	return
}

func (r *repository) ReadRoleAccessCount(
	ctx context.Context,
	args ReadRoleAccessParams,
) (count int, err error) {
	const stmt = `
		SELECT COUNT(*)
		FROM app_menu
		WHERE 
			(CASE WHEN $1::bool THEN (
				name ILIKE $2
			) ELSE TRUE END)
	`
	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)
	return
}
