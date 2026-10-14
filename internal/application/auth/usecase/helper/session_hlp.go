package helper

import (
	"context"

	"be-dashboard-nba/internal/infrastructure/logger"

	"be-dashboard-nba/internal/application/jwt"
	"be-dashboard-nba/internal/domain/model"
	authPresenter "be-dashboard-nba/internal/presentation/auth/presenter"
)

func GenerateSessionModel(
	ctx context.Context,
	request authPresenter.SessionPayload,
) (data model.Session, err error) {
	accessToken, err := jwt.GenerateAccessToken(request.ToAccessTokenRequest())
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to generate access token for user %s: %v", request.UserID, err)
		return
	}

	refreshToken, err := jwt.GenerateRefreshToken(request.ToRefreshTokenRequest())
	if err != nil {
		logger.WithContext(ctx).Errorf("failed to generate refresh token for user %s: %v", request.UserID, err)
		return
	}

	data = model.Session{
		ID:                    request.SessionID,
		UserID:                request.UserID,
		AccessToken:           accessToken.Token,
		AccessTokenExpiredAt:  accessToken.ExpiresAt,
		RefreshToken:          refreshToken.Token,
		RefreshTokenExpiredAt: refreshToken.ExpiresAt,
		IPAddress:             request.IPAddress,
		UserAgent:             request.UserAgent,
	}

	return
}

