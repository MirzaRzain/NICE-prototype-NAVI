package auth

import "github.com/gofiber/fiber/v2"

func GetSessionID(c *fiber.Ctx) (string, error) {
	sessionID, ok := c.Locals("session_id").(string)
	if !ok || sessionID == "" {
		return "", fiber.ErrUnauthorized
	}

	return sessionID, nil
}
