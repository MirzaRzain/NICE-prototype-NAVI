package jwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshTokenPayload struct {
	SessionID string
}

// Generate refresh token with signing JWT method HS256.
func GenerateRefreshToken(request RefreshTokenPayload) (response TokenPayload, err error) {
	claims := &jwt.MapClaims{
		"jti": request.SessionID,
	}
	permanent, _ := strconv.ParseBool(os.Getenv("AUTH_PERMANENT_TOKENS"))
	expiresAt := time.Now()
	if permanent {
		expiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	} else {
		expiredDuration, parseErr := time.ParseDuration(os.Getenv("AUTH_REFRESH_TOKEN_EXPIRES"))
		if parseErr != nil {
			return response, parseErr
		}
		expiresAt = time.Now().Add(expiredDuration)
		(*claims)["exp"] = expiresAt.Unix()
	}

	token, err := GenerateJWT(claims, os.Getenv("AUTH_REFRESH_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = TokenPayload{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return
}

// Parse JWT and return refresh token payload.
func ClaimsRefreshToken(token string) (response RefreshTokenPayload, err error) {
	claims, err := ClaimsJWT(token, os.Getenv("AUTH_REFRESH_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = RefreshTokenPayload{
		SessionID: claims["jti"].(string),
	}

	return
}
