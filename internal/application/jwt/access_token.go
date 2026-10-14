package jwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: implement role when create role usecase

type AccessTokenPayload struct {
	SessionID string
	UserID    string
	// Role   string
}

// Generate access token with signing JWT method HS256.
func GenerateAccessToken(request AccessTokenPayload) (response TokenPayload, err error) {
	claims := &jwt.MapClaims{
		"jti": request.SessionID,
		"uri": request.UserID,
	}
	permanent, _ := strconv.ParseBool(os.Getenv("AUTH_PERMANENT_TOKENS"))
	expiresAt := time.Now()
	if permanent {
		expiresAt = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	} else {
		expiredDuration, parseErr := time.ParseDuration(os.Getenv("AUTH_ACCESS_TOKEN_EXPIRES"))
		if parseErr != nil {
			return response, parseErr
		}
		expiresAt = time.Now().Add(expiredDuration)
		(*claims)["exp"] = expiresAt.Unix()
	}

	token, err := GenerateJWT(claims, os.Getenv("AUTH_ACCESS_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = TokenPayload{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return
}

// Parse JWT and return access token payload.
func ClaimsAccessToken(token string) (response AccessTokenPayload, err error) {
	claims, err := ClaimsJWT(token, os.Getenv("AUTH_ACCESS_TOKEN_SECRET_KEY"))
	if err != nil {
		return
	}

	response = AccessTokenPayload{
		SessionID: claims["jti"].(string),
		UserID:    claims["uri"].(string),
		// Role:   claims["role"].(string),
	}

	return
}
