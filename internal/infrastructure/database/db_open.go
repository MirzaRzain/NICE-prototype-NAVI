package db

import (
	"database/sql"
)

func openSQL(driver, source string, opt *connectionOption) (db *sql.DB, err error) {
	db, err = sql.Open(driver, source)
	if err != nil {
		return
	}

	db.SetMaxOpenConns(opt.maxOpen)
	db.SetMaxIdleConns(opt.maxIdle)
	db.SetConnMaxIdleTime(opt.maxConnIdleTime)
	db.SetConnMaxLifetime(opt.maxLifetime)

	return
}
