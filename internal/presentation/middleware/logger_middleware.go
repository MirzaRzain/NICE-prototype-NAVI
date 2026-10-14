package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	appLogger "be-dashboard-nba/internal/infrastructure/logger"
)

func LoggerMiddleware(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("X-Request-ID", requestID)
		c.SetUserContext(context.WithValue(c.UserContext(), "request_id", requestID))
		startedAt := time.Now()

		appLogger.WithContext(c.UserContext()).Info().
			Str("event", "http_request").
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("query", c.OriginalURL()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Msg("request received")

		err := c.Next()

		entry := appLogger.WithContext(c.UserContext()).Info().
			Str("event", "http_response").
			Int("status", c.Response().StatusCode()).
			Dur("latency", time.Since(startedAt)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("response_bytes", len(c.Response().Body()))

		if err != nil {
			entry = entry.Err(err)
		}
		entry.Msg("response sent")

		return err
	})
}
