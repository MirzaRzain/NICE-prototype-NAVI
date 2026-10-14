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

func (s *useCase) CreateMenuUseCase(
	ctx context.Context,
	request menuPresenter.CreateMenuRequest,
	userID string,
) (err error) {

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

	params := mapper.ToCreateMenuParams(&request, userID)

	var nextSort int

	if params.ParentID.Valid {
		nextSort, err = q.ReadNextSortForParent(ctx, params.ParentID.Int32)
	} else {
		nextSort, err = q.ReadSortForGroup(ctx, params.Group)
	}

	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error getting next sort value for menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params.Sort = nextSort

	err = q.CreateMenuQuery(ctx, params)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Interface("request_payload", params).Msg("error to create menu")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
