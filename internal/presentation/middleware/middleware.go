package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func JWTProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:     []byte("jwt_secret"),
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccess,
	})
}

func jwtSuccess(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	id, _ := strconv.Atoi(fmt.Sprintf("%v", claims["id"]))
	email := fmt.Sprintf("%v", claims["email"])

	c.Locals("user_id", id)
	c.Locals("user_email", email)
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}
