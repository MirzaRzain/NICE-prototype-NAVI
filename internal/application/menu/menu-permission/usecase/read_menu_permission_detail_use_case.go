package menu_permission

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) ReadMenuPermissionDetailUseCase(ctx context.Context, MenuPermissionID int) (data model.MenuPermission, err error) {
	r := s.newMenuPermissionRepo(s.db)

	data, err = r.ReadMenuPermissionByIdQuery(ctx, MenuPermissionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("id", MenuPermissionID).Msg("menu permission detail not found")
			err = constant.ErrMenuPermissionIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", MenuPermissionID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
