package auth

import (
	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth/dto"
	"be-dashboard-nba/internal/application/jwt"
	"context"

	"github.com/pkg/errors"
)

func (s *useCase) LogoutUseCase(
	ctx context.Context,
	claims *jwt.AccessTokenPayload,
	iPAddress string,
) (err error) {
	q := s.newAuthRepo(s.db)
	logRepo := s.newAuthRepo(s.db)

	err = q.DeleteSessionQuery(ctx, claims.SessionID)
	if err != nil {
		log(ctx).Error().Err(err).Str("session_id", claims.SessionID).Msg("error delete session")
		err = errors.WithStack(constant.ErrUnknownSource)
		return
	}

	logRecordArgs := dto.LoginRecord{
		UserID:      claims.UserID,
		AccessToken: claims.SessionID,
		Status:      "logout",
		IPAddress:   iPAddress,
		Type:        "web",
	}
	if logErr := logRepo.CreateLoginRecord(ctx, logRecordArgs); logErr != nil {
		log(ctx).Error().Err(logErr).Msg("failed to create logout record log")
	}

	return nil
}
