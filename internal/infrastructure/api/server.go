package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"be-dashboard-nba/docs"
	"be-dashboard-nba/internal/infrastructure/api/routes"
	runtime "be-dashboard-nba/internal/infrastructure/runtime"
	"be-dashboard-nba/internal/infrastructure/runtime/container"
	"be-dashboard-nba/internal/presentation/middleware"
)

func RunFiberServer(ctx context.Context, c *container.Container, cfg *runtime.Config) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler(),
	})

	middleware.TimeoutMiddleware(app)
	middleware.CorsMiddleware(app)
	middleware.RecoverMiddleware(app)
	middleware.LoggerMiddleware(app)
	middleware.RateLimiterMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": cfg.Name + " is Running",
		})
	})

	if cfg.Swagger.IsEnabled {
		docs.SwaggerInfo.BasePath = cfg.Swagger.BasePath
		docs.SwaggerInfo.Host = cfg.Swagger.Host
		app.Get("/docs/*", swagger.HandlerDefault)
	}

	routes.Routes(app, c, cfg.BaseURLCV)

	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- app.Listen(cfg.Address())
	}()

	c.GetLog().Info().Str("address", cfg.Address()).Msg("starting REST HTTP server")

	select {
	case err := <-listenErrCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.GetLog().Error().Err(err).Msg("server stopped with an error")
			return err
		}
	case <-ctx.Done():
		c.GetLog().Warn().Msg("shutdown signal received, starting graceful shutdown")
	}

	if err := app.ShutdownWithTimeout(cfg.ShutdownTimeoutDuration); err != nil {
		c.GetLog().Error().Err(err).Msg("error shutting down Fiber")
		return err
	}

	time.Sleep(cfg.ShutdownWaitDuration)

	if db := c.GetDB(); db != nil {
		if err := db.Close(); err != nil {
			c.GetLog().Error().Err(err).Msg("error closing database")
			return err
		}
	}

	c.GetLog().Info().Msg("server exited gracefully")
	return nil
}
