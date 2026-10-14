package auth

import (
	"be-dashboard-nba/constant"
	"context"
)

func (s *useCase) CheckPermissionUseCase(ctx context.Context, menuURL constant.MenuKey, userID string, permissionCode []string) (hasAccess bool, err error) {
	r := s.newAuthRepo(s.db)
	hasAccess, err = r.CheckPermissionQuery(ctx, menuURL, userID, permissionCode)

	return
}
