package menu_permission

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/menu-permission/mapper"
	menuPermissionPresenter "be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) UpdateMenuPermissionUseCase(ctx context.Context, payload menuPermissionPresenter.UpdateMenuPermissionRequest, userID string, menuPermissionID int) (err error) {
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
			log(ctx).Warn().Int("id", menuPermissionID).Msg("menu permission detail not found for update")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", menuPermissionID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = r.UpdateMenuPermissionQuery(ctx, mapper.ToUpdateMenuPermissionParams(&payload, userID, menuPermissionID))
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", payload).Msg("error to update menu permission")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
