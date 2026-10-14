package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/role/mapper"
	rolePresenter "be-dashboard-nba/internal/presentation/role/presenter"
	"context"

	"github.com/pkg/errors"
)

func (s *useCase) CreateRoleUseCase(ctx context.Context, paylaod rolePresenter.CreateRoleRequest, userID string) (err error) {
	r := s.newRoleRepo(s.db)

	err = r.CreateRoleQuery(ctx, mapper.ToCreateRoleParams(&paylaod, userID))
	if err != nil {
		log(ctx).Error().Err(err).Interface("request_payload", paylaod).Msg("error to create role")
		err = errors.WithStack(constant.ErrUnknownSource)
	}
	return err
}
