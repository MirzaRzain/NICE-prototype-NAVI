package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"be-dashboard-nba/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
	"be-dashboard-nba/internal/application/jwt"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	db "be-dashboard-nba/internal/infrastructure/database"
	auth "be-dashboard-nba/internal/presentation/auth"
)

type EnsureToken struct {
	db db.DB
}

func NewEnsureToken(db db.DB) *EnsureToken {
	return &EnsureToken{db: db}
}

func (et *EnsureToken) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ah := auth.NewAuth(et.db)

		tokenHeader := c.Get(constant.DefaultMdwHeaderToken)
		token, err := parseHeaderToken(tokenHeader)
		if err != nil {
			logger.WithContext(c.Context()).Error(err, "error parse header token")
			return fiber.NewError(http.StatusUnauthorized, "Token tidak valid")
		}

		accessTokenClaims, err := jwt.ClaimsAccessToken(token)
		if err != nil {
			logger.WithContext(c.Context()).Error(err, "error claims access token")
			return fiber.NewError(http.StatusUnauthorized, "Token tidak valid atau kedaluwarsa")
		}

		ah.SetClaims(&accessTokenClaims)

		err = ah.ValidateSession(c.UserContext())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) ||
				errors.Is(err, constant.ErrTokenExpired) ||
				errors.Is(err, constant.ErrTokenInvalid) ||
				errors.Is(err, constant.ErrTokenUnauthorized) {
				logger.WithContext(c.Context()).Warn(err, "Invalid session (token valid, but session not found)")
				return fiber.NewError(http.StatusUnauthorized, "Sesi tidak valid atau telah berakhir")
			}

			logger.WithContext(c.Context()).Error(err, "Failed to validate session (database error)")
			return fiber.NewError(http.StatusInternalServerError, "Gagal memvalidasi sesi")
		}

		c.Locals("auth", *ah)
		return c.Next()
	}
}

func parseHeaderToken(headerDataToken string) (string, error) {
	if !strings.Contains(headerDataToken, "Bearer") {
		return "", constant.ErrHeaderTokenNotFound
	}

	splitToken := strings.Split(headerDataToken, fmt.Sprintf("%s ", constant.DefaultMdwHeaderBearer))
	if len(splitToken) <= 1 {
		return "", constant.ErrHeaderTokenInvalid
	}

	return splitToken[1], nil
}

func Authorize(
	svc contract.AuthUseCase,
	menuURL constant.MenuKey,
	permissionCode constant.PermissionCode,
) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := auth.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		userID := ah.GetClaims().UserID
		codesToCheck := constant.GetInheritedPermissions(permissionCode)

		hasAccess, err := svc.CheckPermissionUseCase(c.UserContext(), menuURL, userID, codesToCheck)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "Failed to check permissions")
			return fiber.NewError(http.StatusInternalServerError, "Failed to check permissions")

		}
		if !hasAccess {
			return fiber.NewError(http.StatusForbidden, "Access Forbidden")
		}
		return c.Next()
	}
}

