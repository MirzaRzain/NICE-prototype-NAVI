package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/user/mapper"
	"be-dashboard-nba/internal/application/utils"
	userPresenter "be-dashboard-nba/internal/presentation/user/presenter"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) CreateUserUseCase(ctx context.Context, req userPresenter.CreateUserRequest, userID string) (err error) {
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

	ur := s.newUserRepo(tx)
	rr := s.newRoleRepo(tx)

	hashedPassword, err := utils.GenerateHashPassword(req.Password)
	if err != nil {
		log(ctx).Error().Err(err).Msg("failed hash password")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := mapper.ToCreateUserParams(&req, userID, hashedPassword)
	roleID := params.RoleID

	_, err = rr.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("Id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", roleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	newUserID, err := ur.CreateUserQuery(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", req).Msg("error to create user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	err = ur.CreateUserRoleQuery(ctx, params.RoleID, newUserID)
	if err != nil {
		log(ctx).Error().Err(err).Msg("error to create user role")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
