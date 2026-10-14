package menu

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/menu/mapper"
	"be-dashboard-nba/internal/infrastructure/logger"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) UpdateMenuOrderUseCase(ctx context.Context, request menuPresenter.UpdateMenuOrderRequest, userID string) (err error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error to begin transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				logger.WithContext(ctx).Error().Err(errRollback).AnErr("original_error", err).Msg("error to rollback transaction")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}()

	q := s.newMenuRepo(tx)

	paramsList := mapper.ToUpdateMenuSortParamsList(&request, userID)

	for _, params := range paramsList {

		err = q.UpdateMenuSortQuery(ctx, params)
		if err != nil {
			logger.WithContext(ctx).Error().Err(err).Int("menu_id", params.ID).Msg("Failed to update menu sort")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
		err = q.UpdateChildrenGroup(ctx, params.ID, params.Group)
		if err != nil {
			logger.WithContext(ctx).Error().Err(err).Int("menu_id", params.ID).Msg("Failed to update children group on update sort")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
