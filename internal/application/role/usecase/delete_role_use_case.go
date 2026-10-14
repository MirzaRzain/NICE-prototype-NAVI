package role

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) DeleteRoleUseCase(ctx context.Context, roleID int) (err error) {

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log(ctx).Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	r := s.newRoleRepo(tx)

	_, err = r.ReadRoleByIDQuery(ctx, roleID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", roleID).Msg("menu detail not found for update")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", roleID).Msg("error reading menu detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.DeleteRoleQuery(ctx, roleID)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to delete role")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return
}
