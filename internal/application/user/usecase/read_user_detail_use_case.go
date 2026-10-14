package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

func (s *useCase) ReadDetailUserUseCase(
	ctx context.Context,
	id string,
) (data model.User, err error) {
	r := s.newUserRepo(s.db)

	data, err = r.ReadUserByIDQuery(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Str("id", id).Msg("user detail not found")
			err = constant.ErrUserIdNotFound
			return
		}
		log(ctx).Error().Err(err).Str("id", id).Msg("error reading user detail query")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	return
}
