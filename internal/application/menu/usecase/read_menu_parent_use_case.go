package menu

import (
	"be-dashboard-nba/internal/domain/model"
	"be-dashboard-nba/internal/infrastructure/logger"
	"context"
)

func (s *useCase) ReadMenuParentUseCase(ctx context.Context) (data []model.Menu, err error) {
	r := s.newMenuRepo(s.db)
	data, err = r.ReadParentMenuQuery(ctx)
	if err != nil {
		logger.WithContext(ctx).Error().Err(err).Msg("error query read  menu parent")
		return
	}

	return
}
