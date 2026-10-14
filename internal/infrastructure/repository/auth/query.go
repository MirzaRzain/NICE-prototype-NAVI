package repository

import (
	"context"

	"be-dashboard-nba/internal/application/auth/dto"
)

func (r *repository) CreateSessionQuery(
	ctx context.Context,
	args dto.SessionParams,
) (err error) {
	const statement = `
		INSERT INTO sessions (
			id, user_id,
			access_token, access_token_expired_at,
			refresh_token, refresh_token_expired_at,
			ip_address, user_agent,
			created_at
		)
		VALUES (
			$1, $2,
			$3, $4,
			$5, $6,
			$7, $8,
			CURRENT_TIMESTAMP
		)
	`

	_, err = r.DB.ExecContext(ctx, statement,
		args.ID,
		args.UserID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
		args.IPAddress,
		args.UserAgent,
	)

	return
}

func (r *repository) UpdateSessionQuery(
	ctx context.Context,
	args dto.SessionParams,
) (err error) {
	const statement = `
		UPDATE
			sessions
		SET
			access_token = $2,
			access_token_expired_at = $3,
			refresh_token = $4,
			refresh_token_expired_at = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $1
	`

	_, err = r.DB.ExecContext(ctx, statement,
		args.ID,
		args.AccessToken,
		args.AccessTokenExpiredAt,
		args.RefreshToken,
		args.RefreshTokenExpiredAt,
	)

	return
}

func (r *repository) DeleteSessionQuery(
	ctx context.Context,
	id string,
) (err error) {

	const statement = `
		DELETE FROM sessions
		WHERE
			id = $1
	`

	_, err = r.DB.ExecContext(ctx, statement, id)

	return
}

type LoginAttempParams = dto.LoginAttempParams

func (r *repository) CreateLoginAttemp(
	ctx context.Context,
	args LoginAttempParams,
) (err error) {

	const statement = `
		INSERT INTO app_login_attempt (
			pin, ip_address
		)
		VALUES ($1, $2)
	`
	_, err = r.DB.ExecContext(ctx, statement,
		args.Password,
		args.IPAddress,
	)
	return
}

type LoginRecord = dto.LoginRecord

func (r *repository) CreateLoginRecord(
	ctx context.Context,
	args LoginRecord,
) (err error) {

	const statement = `
		INSERT INTO app_login (
			access_token, status, ip_address, type
		)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.DB.ExecContext(ctx, statement,
		args.AccessToken,
		args.Status,
		args.IPAddress,
		args.Type,
	)
	return
}
