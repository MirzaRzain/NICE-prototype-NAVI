package routes

import (
	"github.com/gofiber/fiber/v2"

	usecase "be-dashboard-nba/internal/application/user/usecase"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	"be-dashboard-nba/internal/presentation/middleware"
	handlers "be-dashboard-nba/internal/presentation/user/handler"
)

func UserRouter(http fiber.Router, c *container.Container) {
	svc := usecase.NewUseCase(c.GetDB())
	mdw := middleware.NewEnsureToken(c.GetDB())

	routes := http.Group("/users")
	routes.Use(mdw.ValidateToken())

	routes.Get("/", handlers.ReadUsers(svc))
	routes.Post("/", handlers.CreateUser(svc, c.GetValidator()))
	routes.Put("/:user_id", handlers.UpdateUser(svc, c.GetValidator()))
	routes.Delete("/:user_id", handlers.DeleteUser(svc))
	routes.Put("/me", handlers.UpdateProfileApp(svc, c.GetValidator()))
	routes.Get("/me", handlers.ReadProfileApp(svc))
	routes.Get("/:user_id", handlers.ReadUserDetail(svc))

}
