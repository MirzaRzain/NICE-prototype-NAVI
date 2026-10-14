package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"be-dashboard-nba/constant"
)

func CorsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodHead,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodConnect,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: "Origin, Accept, Content-Type, " + constant.DefaultMdwHeaderToken,
	}))
}
