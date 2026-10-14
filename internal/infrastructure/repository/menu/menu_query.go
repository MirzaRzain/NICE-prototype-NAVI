package repository

import (
	"context"

	"be-dashboard-nba/internal/application/menu/dto"
)

func (r *repository) CreateMenuQuery(
	ctx context.Context,
	params dto.CreateMenuParams,
) (err error) {
	const stmt = `
		INSERT INTO app_menu (
			parent_id, name, description, url, sort, "group", icon,
			active, display, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)
		
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.ParentID, params.Name, params.Description, params.URL, params.Sort,
		params.Group, params.Icon, params.Active, params.Display, params.CreatedBy,
	)
	return
}

func (r *repository) UpdateMenuQuery(
	ctx context.Context,
	params dto.UpdateMenuParams,
) (err error) {
	const stmt = `
        UPDATE app_menu
        SET
            name = $2, 
            parent_id = $3,
            description = $4,
            url = $5,
            sort = $6,
            "group" = $7, 
            icon = $8,      
            active = $9,     
            display = $10, 
            updated_by = $11,
            updated_at = CURRENT_TIMESTAMP
        WHERE
            id = $1;
    `

	_, err = r.db.ExecContext(ctx, stmt,
		params.ID,
		params.Name,
		params.ParentID,
		params.Description,
		params.URL,
		params.Sort,
		params.Group,
		params.Icon,
		params.Active,
		params.Display,
		params.UpdatedBy,
	)
	return
}

func (r *repository) DeleteMenuQuery(ctx context.Context, menuID int) (err error) {
	const stmt = `
	DELETE FROM app_menu
	WHERE 
		id = $1
`

	_, err = r.db.ExecContext(ctx, stmt, menuID)

	return
}

type UpdateMenuSortParams = dto.UpdateMenuSortParams

func (r *repository) UpdateMenuSortQuery(
	ctx context.Context,
	params UpdateMenuSortParams,
) (err error) {
	const stmt = `
		UPDATE app_menu
		SET 
			sort = $1,
			"group" = $2,
			parent_id = $3,
			updated_by = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE 
			id = $5
	`
	_, err = r.db.ExecContext(ctx, stmt,
		params.Sort,
		params.Group,
		params.ParentID,
		params.UpdatedBy,
		params.ID,
	)
	return
}

func (r *repository) UpdateChildrenGroup(ctx context.Context, parentID int, newGroup string) error {
	const stmt = `UPDATE app_menu SET "group" = $1, updated_at = CURRENT_TIMESTAMP WHERE parent_id = $2`
	_, err := r.db.ExecContext(ctx, stmt, newGroup, parentID)
	return err
}
