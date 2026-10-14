package model

import (
	"database/sql"
	"time"
)

type Session struct {
	ID                    string       `db:"id"`      // DIUBAH: dari GUID ke ID (uuid -> string)
	UserID                string       `db:"user_id"` // DIUBAH: dari int64 ke string
	AccessToken           string       `db:"access_token"`
	AccessTokenExpiredAt  time.Time    `db:"access_token_expired_at"`
	RefreshToken          string       `db:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `db:"refresh_token_expired_at"`
	IPAddress             string       `db:"ip_address"`
	UserAgent             string       `db:"user_agent"`
	CreatedAt             sql.NullTime `db:"created_at"` // DIUBAH: DDL baru Anda 'NULL'
	UpdatedAt             sql.NullTime `db:"updated_at"`
}
