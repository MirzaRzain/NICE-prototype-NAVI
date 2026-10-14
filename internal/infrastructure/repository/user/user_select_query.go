package repository

import (
	"context"

	"be-dashboard-nba/internal/application/user/dto"
	"be-dashboard-nba/internal/domain/model"
)

func (r *repository) ReadUsersQuery(
	ctx context.Context,
	args dto.ReadListUserParams,
) (data []model.User, err error) {
	const stmt = `
		SELECT
			au.id, 
			au.full_name, 
			au.active,
			ar.name AS role,
			ar.id AS role_id
		FROM
			app_user au
		JOIN 
			app_user_role aur ON au.id = aur.user_id
		JOIN 
			app_role ar ON aur.role_id = ar.id
		WHERE
			(CASE WHEN $1::bool THEN(
				full_name ILIKE $2
			) ELSE TRUE END)
		ORDER BY
			(CASE WHEN $3 = 'full_name ASC' THEN au.full_name END) ASC,
			(CASE WHEN $3 = 'full_name DESC' THEN au.full_name END) DESC,
			(CASE WHEN $3 = 'role ASC' THEN ar.name END) ASC,
			(CASE WHEN $3 = 'role DESC' THEN ar.name END) DESC,
			(CASE WHEN $3 = 'active ASC' THEN au.active END) ASC,
			(CASE WHEN $3 = 'active DESC' THEN au.active END) DESC,
			(CASE WHEN $3 = 'created_at ASC' THEN au.created_at END) ASC,
			(CASE WHEN $3 = 'created_at DESC' THEN au.created_at END) DESC
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
		var u model.User

		if err = rows.Scan(
			&u.ID,
			&u.FullName,
			&u.Active,
			&u.Role,
			&u.RoleID,
		); err != nil {
			return
		}

		data = append(data, u)
	}

	return
}

func (r *repository) ReadCountUserQuery(
	ctx context.Context,
	args dto.ReadListUserParams,
) (count int, err error) {
	const stmt = `
		SELECT
			COUNT(*)
		FROM
			app_user
		WHERE
			(CASE WHEN $1::bool THEN(
				full_name ILIKE $2
			) ELSE TRUE END)
	`

	err = r.db.QueryRowContext(ctx, stmt,
		args.SetSearch,
		args.Search,
	).Scan(&count)

	return
}

func (r *repository) ReadUserByIDQuery(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	const statement = `
		SELECT 	
			au.id, 
			au.name, 
			au.full_name, 
			au.email,
			au.phone,
			au.active,
			ar.name AS role,
			ar.id AS role_id
		FROM
			app_user au
		JOIN
			app_user_role aur ON au.id = aur.user_id
		JOIN 
			app_role ar ON aur.role_id = ar.id
		WHERE
			au.id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Email,
		&data.Phone,
		&data.Active,
		&data.Role,
		&data.RoleID,
	)

	return
}

func (r *repository) ReadUserProfileQuery(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	const statement = `
		SELECT 	
			au.id, 
			au.name, 
			au.full_name, 
			au.email,
			au.phone,
			au.active,
			au.img_path,
			au.img_name,
			ar.name AS role,
			ar.id AS role_id
		FROM
			app_user au
		JOIN
			app_user_role aur ON au.id = aur.user_id
		JOIN 
			app_role ar ON aur.role_id = ar.id
		WHERE
			au.id = $1
	`

	err = r.db.QueryRowContext(ctx, statement, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Email,
		&data.Phone,
		&data.Active,
		&data.ImgPath,
		&data.ImgName,
		&data.Role,
		&data.RoleID,
	)

	return
}

func (r *repository) IsUserEmailExistsQuery(
	ctx context.Context,
	email string,
) (exists bool, err error) {
	const statement = `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email).Scan(&exists)

	return
}

func (r *repository) IsUpdateUserEmailExistsQuery(
	ctx context.Context,
	email, id string,
) (exists bool, err error) {
	statement := `
		SELECT EXISTS (
			SELECT
				1
			FROM
				app_user
			WHERE
				email = $1
				AND id != $2
		)
	`

	err = r.db.QueryRowContext(ctx, statement, email, id).Scan(&exists)

	return
}
