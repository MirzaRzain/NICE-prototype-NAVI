package repository

import (
	"be-dashboard-nba/internal/application/menu/dto"
	"be-dashboard-nba/internal/domain/model"
	"context"
)

func (r *repository) ReadSidebarMenuQuery(
	ctx context.Context,
	userID string,
) (data []model.Menu, err error) {

	const stmt = `
		SELECT DISTINCT
			m.id, m.parent_id, m.name, m.url, m.sort,
			m."group", m.icon, m.active, m.display
		FROM
			app_menu m
		JOIN
			app_menu_permission amp ON m.id = amp.menu_id
		JOIN
			app_role_access ara ON amp.id = ara.menu_permission_id
		JOIN
			app_user_role aur ON ara.role_id = aur.role_id
		WHERE
			aur.user_id = $1
			AND amp.code = 'R' 
			AND m.display = true   
			AND m.active = true  
		ORDER BY
    		m."group" ASC,
   		 	m.sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}

func (r *repository) ReadListMenuQuery(
	ctx context.Context,
	params dto.ReadListMenuParams,
) (data []model.Menu, err error) {
	const stmt = `
		WITH RECURSIVE menu_with_parents AS (
			SELECT
				id, parent_id, name, description, url, sort, "group", icon,
				active, display, created_by, created_at, updated_by, updated_at
			FROM
				app_menu
			WHERE
				(CASE WHEN $1::bool THEN
					name ILIKE $2
				ELSE TRUE 
				END)

			UNION

			SELECT
				m.id, m.parent_id, m.name, m.description, m.url, m.sort, m."group", m.icon,
				m.active, m.display, m.created_by, m.created_at, m.updated_by, m.updated_at
			FROM
				app_menu m
			JOIN
				menu_with_parents mwp ON m.id = mwp.parent_id 
		)
		SELECT DISTINCT * FROM menu_with_parents
		ORDER BY
			sort ASC;
	`

	rows, err := r.db.QueryContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.ParentID, &m.Name, &m.Description, &m.URL, &m.Sort,
			&m.Group, &m.Icon, &m.Active, &m.Display,
			&m.CreatedBy, &m.CreatedAt, &m.UpdatedBy, &m.UpdatedAt,
		); err != nil {
			return
		}
		data = append(data, m)
	}

	return
}

func (r *repository) ReadCountMenuQuery(
	ctx context.Context,
	params dto.ReadListMenuParams,
) (count int64, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_menu
		WHERE
			(CASE WHEN $1::bool THEN(
				name ILIKE $2
				OR description ILIKE $2
				OR url ILIKE $2
			) ELSE TRUE END)
	`
	err = r.db.QueryRowContext(ctx, stmt,
		params.SetSearch,
		params.Search,
	).Scan(&count)
	return
}

func (r *repository) ReadParentMenuQuery(
	ctx context.Context,
) (data []model.Menu, err error) {
	const stmt = `
		SELECT
			id, name, "group"
		FROM
			app_menu
		WHERE
			parent_id IS NULL
		ORDER BY
			sort ASC;
	`
	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m model.Menu
		if err = rows.Scan(
			&m.ID, &m.Name, &m.Group,
		); err != nil {
			return
		}
		data = append(data, m)
	}
	return
}

func (r *repository) ReadMenuByIDQuery(ctx context.Context, menuID int) (data model.Menu, err error) {
	const stmt = `
		SELECT
			id, parent_id, name, description, url, sort, "group", icon,
			active, display
		FROM 
			app_menu
		WHERE
			id = $1
	`
	err = r.db.QueryRowContext(ctx, stmt, menuID).Scan(
		&data.ID, &data.ParentID, &data.Name, &data.Description, &data.URL, &data.Sort,
		&data.Group, &data.Icon, &data.Active, &data.Display,
	)
	return
}

func (r *repository) ReadNextSortForParent(ctx context.Context, parentID int32) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1 
		FROM app_menu
		WHERE parent_id = $1
	`

	err := r.db.QueryRowContext(ctx, stmt, parentID).Scan(&nextSort)
	return nextSort, err
}

func (r *repository) ReadSortForGroup(ctx context.Context, group string) (int, error) {
	var nextSort int
	const stmt = `
		SELECT COALESCE(MAX(sort), -1) + 1 
		FROM app_menu
		WHERE parent_id IS NULL AND "group" = $1
	`
	err := r.db.QueryRowContext(ctx, stmt, group).Scan(&nextSort)
	return nextSort, err
}

func (r *repository) ReadNextSortForParentAndGroup(ctx context.Context, parentID int32, group string) (int, error) {
	const stmt = `SELECT COALESCE(MAX(sort), -1) + 1 FROM app_menu WHERE parent_id = $1 AND "group" = $2`
	var nextSort int
	err := r.db.QueryRowContext(ctx, stmt, parentID, group).Scan(&nextSort)
	return nextSort, err
}

func (r *repository) CountMenuChildren(ctx context.Context, menuID int) (int, error) {
	const stmt = `SELECT COUNT(*) FROM app_menu WHERE parent_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, stmt, menuID).Scan(&count)
	return count, err
}
