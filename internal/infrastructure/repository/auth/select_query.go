package repository

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"

	"github.com/lib/pq"
)

func (r *repository) ReadDetailUserByEmailQuery(
	ctx context.Context,
	email string,
) (data model.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			email = $1
	`
	err = r.DB.QueryRowContext(ctx, statement, email).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Password,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}

func (r *repository) ReadDetailUserByIdQuery(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	const statement = `
		SELECT
			id, name, full_name, password ,email, active, phone, img_path, img_name
		FROM
			app_user
		WHERE
			id = $1
	`
	err = r.DB.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Password,
		&data.Email,
		&data.Active,
		&data.Phone,
		&data.ImgPath,
		&data.ImgName,
	)

	return
}

func (r *repository) ReadDetailSessionQuery(
	ctx context.Context,
	id string,
) (data model.Session, err error) {
	const statement = `
		SELECT
			id, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at
		FROM
			sessions
		WHERE
			id = $1
	`

	err = r.DB.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.UserID,
		&data.AccessToken,
		&data.AccessTokenExpiredAt,
		&data.RefreshToken,
		&data.RefreshTokenExpiredAt,
	)

	return
}

func (r *repository) CheckPermissionQuery(
	ctx context.Context,
	menuURL constant.MenuKey,
	userID string,
	permissionCodes []string,
) (bool, error) {

	const stmt = `
		WITH RECURSIVE menu_ancestry AS (
    		SELECT id, parent_id, url
    		FROM app_menu
    		WHERE url = $1

    		UNION ALL
    
    		SELECT m.id, m.parent_id, m.url
    		FROM app_menu m
    		JOIN menu_ancestry ma ON m.id = ma.parent_id
		),

		child_menu AS (
    		SELECT id FROM menu_ancestry WHERE url = $1
		),

		parent_menus AS (
    		SELECT id FROM menu_ancestry WHERE url IS DISTINCT FROM $1
		),

		permitted_parents AS (
    		SELECT DISTINCT ma.id
    		FROM parent_menus ma
    		JOIN app_menu_permission amp ON ma.id = amp.menu_id
    		JOIN app_role_access ara ON amp.id = ara.menu_permission_id
    		JOIN app_user_role aur ON ara.role_id = aur.role_id
    		WHERE aur.user_id = $2
      		AND amp.code = ANY($3)
		),

		child_permission_check AS (
    		SELECT 
        	COUNT(DISTINCT amp.code) AS total_permissions
    		FROM child_menu cm
    		JOIN app_menu_permission amp ON cm.id = amp.menu_id
    		JOIN app_role_access ara ON amp.id = ara.menu_permission_id
    		JOIN app_user_role aur ON ara.role_id = aur.role_id
    		WHERE 
        		aur.user_id = $2
        		AND amp.code = ANY($3)
		)

		SELECT (
        	(SELECT COUNT(*) FROM parent_menus) = 0 
        	OR 
        	(SELECT COUNT(*) FROM permitted_parents) > 0
    	)
    	AND
	    (
        (SELECT total_permissions FROM child_permission_check) = cardinality($3)
   		 ) AS has_access;
	`

	var hasAccess bool
	err := r.DB.QueryRowContext(
		ctx,
		stmt,
		menuURL.String(),
		userID,
		pq.Array(permissionCodes),
	).Scan(&hasAccess)

	if err != nil {
		return false, err
	}

	return hasAccess, nil
}
