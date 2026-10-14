package user

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/domain/model"
	"context"
	"database/sql"
	"errors"
)

func (s *useCase) ReadUserProfileUseCase(ctx context.Context, userID string) (data model.User, err error) {

	r := s.newUserRepo(s.db)

	data, err = r.ReadUserProfileQuery(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log(ctx).Warn().Str("id", userID).Msg("user profile not found")
			err = constant.ErrUserIdNotFound
			return
		}
		log(ctx).Error().Err(err).Str("id", userID).Msg("error reading query user pprofile")
		err = constant.ErrUnknownSource
		return
	}

	return
}
