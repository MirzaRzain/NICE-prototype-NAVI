package auth

import (
	"context"

	appLogger "be-dashboard-nba/internal/infrastructure/logger"
)

func log(ctx context.Context) *appLogger.Logger {
	return appLogger.WithContext(ctx)
}
