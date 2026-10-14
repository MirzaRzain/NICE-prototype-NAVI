package auth

import (
	"context"
	"time"

	"be-dashboard-nba/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/jwt"
	db "be-dashboard-nba/internal/infrastructure/database"
	repository "be-dashboard-nba/internal/infrastructure/repository/auth"
)

type Auth struct {
	db     db.DB
	claims *jwt.AccessTokenPayload
}

func NewAuth(db db.DB) *Auth {
	return &Auth{
		db: db,
	}
}

func GetAuth(c *fiber.Ctx) (*Auth, error) {
	a, ok := c.Locals("auth").(Auth)
	if !ok {
		return nil, constant.ErrTokenUnauthorized
	}

	return &a, nil
}

func (a *Auth) GetClaims() *jwt.AccessTokenPayload {
	return a.claims
}

func (a *Auth) SetClaims(claims *jwt.AccessTokenPayload) {
	a.claims = claims
}

func (a *Auth) ValidateSession(ctx context.Context) (err error) {
	r := repository.NewRepo(a.db)

	session, err := r.ReadDetailSessionQuery(ctx, a.claims.SessionID)
	if err != nil {
		logger.WithContext(ctx).Error(err, "error read session by id : "+a.claims.SessionID)
		return
	}

	if time.Now().After(session.AccessTokenExpiredAt) {
		err = constant.ErrTokenExpired
		return
	}

	return
}

