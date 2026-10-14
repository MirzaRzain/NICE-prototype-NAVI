package middleware

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"be-dashboard-nba/constant"
)

func TimeoutMiddleware(app *fiber.App) {
	timeoutStr := os.Getenv("APP_REQUEST_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil || timeout <= 0 {
		timeout = constant.DefaultMdwTimeout
	}

	app.Use(func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		c.SetUserContext(ctx)

		err := c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			return fiber.NewError(fiber.StatusRequestTimeout, "Request Timeout")
		}

		return err
	})
}
