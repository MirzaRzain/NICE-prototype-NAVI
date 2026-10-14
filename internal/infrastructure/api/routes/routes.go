package routes

import (
	"be-dashboard-nba/internal/infrastructure/runtime/container"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, c *container.Container, baseURLCV string) {
	api := app.Group("/api").Group("/v1")
	CaptureRouter(api, c, baseURLCV)
	AuthRouter(api, c)
	MasterRouter(api, c)
}
