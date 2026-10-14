package auth

import (
	"context"
	"os"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/auth/dto"
	"be-dashboard-nba/internal/application/auth/usecase/helper"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/auth/presenter"

	"github.com/pkg/errors"
)

const deviceIdentity = "device"

func (s *useCase) LoginUseCase(ctx context.Context, request authPresenter.LoginRequest, userAgent, iPAddress string) (data model.Session, user model.User, err error) {
	q := s.newAuthRepo(s.db)
	if request.PIN == "" || request.PIN != os.Getenv("AUTH_PIN") {
		if logErr := q.CreateLoginAttemp(ctx, dto.LoginAttempParams{Email: "device", Password: "", IPAddress: iPAddress}); logErr != nil {
			log(ctx).Error().Err(logErr).Msg("failed to create login attempt log")
		}
		return data, user, constant.ErrWrongEmailOrPassword
	}

	data, err = helper.GenerateSessionModel(ctx, request.ToSessionPayload(deviceIdentity, userAgent, iPAddress))
	if err != nil {
		return data, user, errors.WithStack(constant.ErrUnknownSource)
	}

	if err = q.CreateSessionQuery(ctx, dto.SessionParams{ID: data.ID, UserID: data.UserID, AccessToken: data.AccessToken, AccessTokenExpiredAt: data.AccessTokenExpiredAt, RefreshToken: data.RefreshToken, RefreshTokenExpiredAt: data.RefreshTokenExpiredAt, IPAddress: data.IPAddress, UserAgent: data.UserAgent}); err != nil {
		return data, user, errors.WithStack(constant.ErrUnknownSource)
	}
	if logErr := q.CreateLoginRecord(ctx, dto.LoginRecord{UserID: deviceIdentity, AccessToken: data.AccessToken, Status: "login", IPAddress: iPAddress, Type: "api"}); logErr != nil {
		log(ctx).Error().Err(logErr).Msg("failed to create login record log")
	}
	return data, user, nil
}
