package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RecoverMiddleware(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
}
