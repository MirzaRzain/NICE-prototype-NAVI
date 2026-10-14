package menu_permission

import (
	"be-dashboard-nba/constant"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) DeleteMenuPermissionUseCase(ctx context.Context, menuPermissionID int) (err error) {

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

	r := s.newMenuPermissionRepo(tx)

	_, err = r.ReadMenuPermissionByIdQuery(ctx, menuPermissionID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", menuPermissionID).Msg("menu detail not found for delete")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", menuPermissionID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.DeleteMenuPermissionQuery(ctx, menuPermissionID)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to delete menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
