package routes

import (
	masterUseCase "be-dashboard-nba/internal/application/master/usecase"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	masterHandler "be-dashboard-nba/internal/presentation/master/handler"

	"github.com/gofiber/fiber/v2"
)

func MasterRouter(http fiber.Router, c *container.Container) {
	svc := masterUseCase.NewUseCase(c.GetDB())
	routes := http.Group("/master")

	routes.Get("/usecases", masterHandler.ReadUsecases(svc))
	routes.Get("/usecases/config", masterHandler.ReadUsecasesWithConfig(svc))
	routes.Post("/usecases", masterHandler.CreateUsecase(svc, c.GetValidator()))
	routes.Get("/usecases/:usecase_id", masterHandler.ReadUsecase(svc))
	routes.Put("/usecases/:usecase_id", masterHandler.UpdateUsecase(svc, c.GetValidator()))
	routes.Post("/usecases/:usecase_id/inference-objects", masterHandler.SetInferenceObjects(svc, c.GetValidator()))
	routes.Get("/usecases/:usecase_id/inference-objects", masterHandler.ReadInferenceObjects(svc))
	routes.Put("/usecases/:usecase_id/inference-objects", masterHandler.SetInferenceObjects(svc, c.GetValidator()))

	routes.Get("/models", masterHandler.ReadModels(svc))
	routes.Post("/models", masterHandler.CreateModel(svc, c.GetValidator()))
	routes.Get("/models/download/:usecase_code", masterHandler.DownloadModel(svc))
	routes.Get("/models/:model_id", masterHandler.ReadModel(svc))
	routes.Put("/models/:model_id", masterHandler.UpdateModel(svc, c.GetValidator()))
	routes.Delete("/models/:model_id", masterHandler.DeleteModel(svc))
}
