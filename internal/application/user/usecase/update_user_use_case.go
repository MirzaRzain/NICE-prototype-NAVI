package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/user/mapper"
	presenter "be-dashboard-nba/internal/presentation/user/presenter"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) UpdateUserUseCase(
	ctx context.Context,
	request presenter.UpdateUserRequest,
	updatedBy string,
	userID string,
) (err error) {
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

	ru := s.newUserRepo(tx)
	rr := s.newRoleRepo(tx)

	existingUser, err := ru.ReadUserByIDQuery(ctx, userID)
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

	existingRoleUser, err := rr.ReadRoleByIDQuery(ctx, request.RoleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", request.RoleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", request.RoleID).Msg("error reading role detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	params := mapper.ToUpdateUserParams(&request, userID, updatedBy, existingUser)

	if existingUser.Role != existingRoleUser.Name {
		err = ru.UpdateUserRoleQuery(ctx, params.RoleID, params.ID)
		if err != nil {
			log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to update user role")
			err = errors.WithStack(constant.ErrUnknownSource)
			return
		}
	}

	err = ru.UpdateUserQuery(ctx, params)
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", request).Msg("error to update user")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	if err = tx.Commit(); err != nil {
		log(ctx).Error().Err(err).Msg("error to commit transaction")
		err = errors.WithStack(constant.ErrUnknownSource)
	}

	return
}
