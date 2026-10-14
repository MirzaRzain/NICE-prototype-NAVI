package routes

import (
	captureApplication "be-dashboard-nba/internal/application/capture"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	captureHandler "be-dashboard-nba/internal/presentation/capture"

	"github.com/gofiber/fiber/v2"
)

func CaptureRouter(http fiber.Router, c *container.Container, baseURLCV string) {
	service := captureApplication.NewService(c.GetDB(), baseURLCV)
	handler := captureHandler.Proxy(service)

	http.Post("/cv/capture", handler)
	http.Post("/capture", handler)
}
