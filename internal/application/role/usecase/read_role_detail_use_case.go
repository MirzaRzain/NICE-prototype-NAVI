package role

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) ReadRoleDetailUseCase(ctx context.Context, roleID int) (data model.Role, err error) {
	r := s.newRoleRepo(s.db)

	data, err = r.ReadRoleByIDQuery(ctx, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Int("Id", roleID).Msg("role detail not found")
			err = constant.ErrRoleIdNotFound
			return
		}
		log(ctx).Error().Err(err).Int("id", roleID).Msg("error reading menu permission detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
