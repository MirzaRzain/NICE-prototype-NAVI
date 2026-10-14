package routes

import (
	"github.com/gofiber/fiber/v2"

	usecase "be-dashboard-nba/internal/application/auth/usecase"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	handlers "be-dashboard-nba/internal/presentation/auth/handler"
	"be-dashboard-nba/internal/presentation/middleware"
)

func AuthRouter(http fiber.Router, c *container.Container) {
	svc := usecase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("/auth")

	routes.Post("/login", handlers.Login(svc, c.GetValidator()))
	routes.Post("/logout", mdw.ValidateToken(), handlers.Logout(svc))
}
