package db

import (
	"context"
	"database/sql"

	"be-dashboard-nba/constant"

	"github.com/pkg/errors"
)

func WithTransaction(ctx context.Context, db DB, fn func(tx Query) error) (err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.WithStack(constant.ErrUnknownSource)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err = fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.WithStack(constant.ErrUnknownSource)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return errors.WithStack(constant.ErrUnknownSource)
	}

	return nil
}
