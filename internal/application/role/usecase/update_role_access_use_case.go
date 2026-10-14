package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/mapper"
	presenter "be-dashboard-nba/internal/presentation/role/presenter"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) UpdateRoleAccessUseCase(ctx context.Context, roleID int, request presenter.UpdateRoleAccessRequest) (err error) {
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

	rolerepo := s.newRoleRepo(tx)
	menuPermissionRepo := s.newMenuPermissionRepo(tx)
	for _, req := range request.AccessItem {
		params := mapper.ToUpdateRoleMenuPermissionParam(&req, roleID)

		_, err = rolerepo.ReadRoleByIDQuery(ctx, params.RoleID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log(ctx).Warn().Int("id", params.RoleID).Msg("role detail not found for update role menu permission")
				err = constant.ErrRoleIdNotFound
				return
			}
			log(ctx).Error().Err(err).Int("id", params.RoleID).Msg("error reading role detail query on update role menu permission")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}

		_, err = menuPermissionRepo.ReadMenuPermissionByIdQuery(ctx, params.MenuPermissionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log(ctx).Warn().Int("id", params.RoleID).Msg("menu permission detail not found for update role menu permission")
				err = constant.ErrMenuPermissionIdNotFound
				return
			}
			log(ctx).Error().Err(err).Int("id", params.RoleID).Msg("error reading menu permission detail query on update role menu permission")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
		if *req.HasAccess {
			err = rolerepo.CreateRoleAccess(ctx, params)
			if err != nil {
				log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to create role menu permission")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		} else {
			err = rolerepo.DeleteRoleAccess(ctx, params)
			if err != nil {
				log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to delete role menu permission")
				err = errors.WithStack(constant.ErrUnknownSource)
				return
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return

}
