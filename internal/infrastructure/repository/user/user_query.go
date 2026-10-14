package repository

import (
	"context"

	"be-dashboard-nba/internal/application/user/dto"
)

func (r *repository) CreateUserQuery(
	ctx context.Context,
	args dto.CreateUserParams,
) (userID string, err error) {
	const statement = `
		INSERT INTO app_user (
			name, full_name, email, password, active, phone, created_by, created_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, (now() at time zone 'UTC')::TIMESTAMP
		)
		RETURNING id
	`

	err = r.db.QueryRowContext(ctx, statement,
		args.Name,
		args.FullName,
		args.Email,
		args.Password,
		args.Active,
		args.Phone,
		args.CreatedBy,
	).Scan(&userID)

	return
}

func (r *repository) CreateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error) {
	const stmt = `
		INSERT INTO app_user_role(
			user_id, role_id
		)
		VALUES (
			$1, $2
		)
	`

	_, err = r.db.ExecContext(ctx, stmt, userID, roleID)

	return
}

func (r *repository) UpdateUserRoleQuery(ctx context.Context, roleID int, userID string) (err error) {
	const stmt = `
		UPDATE app_user_role
		SET
			user_id = $1,
			role_id = $2
		WHERE user_id = $1
	`

	_, err = r.db.ExecContext(ctx, stmt, userID, roleID)

	return
}

func (r *repository) UpdateUserQuery(
	ctx context.Context,
	args dto.UpdateUserParams,
) (err error) {
	const statement = `
		UPDATE app_user 
		SET 
    		name = $1,
    		full_name = $2,
    		email = $3,
    		phone = $4,
    		active = $5,
    		updated_by = $6,
    		updated_at = (now() at time zone 'UTC')::TIMESTAMP
		WHERE id = $7
	`

	_, err = r.db.ExecContext(ctx, statement,
		args.Name,
		args.FullName,
		args.Email,
		args.Phone,
		args.Active,
		args.UpdatedBy,
		args.ID,
	)

	return
}

func (r *repository) DeleteUserQuery(
	ctx context.Context,
	id string,
) (err error) {
	const statement = `
		DELETE FROM app_user
		WHERE
			id = $1
	`

	_, err = r.db.ExecContext(ctx, statement, id)

	return
}
