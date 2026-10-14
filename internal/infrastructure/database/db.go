package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func NewDatabase(ctx context.Context) (db *sql.DB, err error) {
	opt, err := newDatabaseOption(ctx)
	if err != nil {
		return
	}

	if opt.driver != "sqlite" {
		err = errors.Wrapf(errors.New("invalid datasources driver"), "db: driver=%s", opt.driver)
		return
	}
	db, err = NewSQLite(opt)

	return
}
