package menu

import (
	"be-dashboard-nba/internal/application/menu/mapper"
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/infrastructure/logger"
	menuPresenter "be-dashboard-nba/internal/presentation/menu/presenter"
	"context"
)

func (s *useCase) ReadListMenuUseCase(
	ctx context.Context,
	request menuPresenter.ReadMenuListRequest,
) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadListMenuQuery(ctx, mapper.ToReadListMenuParams(&request))
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error query read list menu")
		return
	}

	return
}
