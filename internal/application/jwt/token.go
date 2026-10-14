package jwt

import (
	"be-dashboard-nba/constant"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type TokenPayload struct {
	Token     string
	ExpiresAt time.Time
}

// Generate JWT with signing method HS256.
func GenerateJWT(claims jwt.Claims, secretKey string) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secretKey))

	return
}

// Parse JWT, and return claims object.
func ClaimsJWT(token, secretKey string) (claims jwt.MapClaims, err error) {
	jwtToken, err := jwt.ParseWithClaims(token, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, errors.Wrapf(err, "Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return
	}

	if jwtToken == nil || !jwtToken.Valid {
		err = errors.Wrap(constant.ErrTokenInvalid, "jwt token is nil or invalid")
		return
	}

	return
}

// Parse Unverified JWT, and return claims object.
func ClaimsUnverifiedJWT(token string) (claims jwt.MapClaims, err error) {
	jwtToken, _, err := new(jwt.Parser).ParseUnverified(token, &claims)
	if err != nil {
		return
	}

	if jwtToken == nil {
		err = errors.Wrap(constant.ErrTokenInvalid, "jwt token is nil or invalid")
		return
	}

	return
}
