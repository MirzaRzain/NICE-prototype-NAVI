package user

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func (s *useCase) DeleteUserUseCase(ctx context.Context, userID string, deltedBy string) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	var originalErr error

	defer func() {
		if err != nil {
			// Simpan error asli sebelum rollback
			originalErr = err
			if errRollback := tx.Rollback(); errRollback != nil {
				log(ctx).Error().Err(errRollback).AnErr("original_error", originalErr).Msg("error to rollback transaction")
				// Kembalikan error asli, jangan overwrite
				err = originalErr
				return
			}
			// Jika rollback success, kembalikan error asli
			err = originalErr
		}
	}()
	r := s.newUserRepo(tx)

	existingUser, err := r.ReadUserByIDQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Str("id", userID).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		log(ctx).Error().Err(err).Str("id", userID).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}
	fmt.Printf("user id: %s, deletedBy: %s", existingUser.ID, deltedBy)
	if existingUser.ID == deltedBy {
		err = constant.ErrForbiddenSelfDelete
		return
	}

	err = r.DeleteUserQuery(ctx, userID)
	if err != nil {
		log(ctx).Error().Err(err).Str("id", userID).Msg("error to delete user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return

}
