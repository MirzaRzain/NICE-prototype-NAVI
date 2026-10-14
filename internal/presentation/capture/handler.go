package capture

import (
	"errors"
	"net/http"

	captureApplication "be-dashboard-nba/internal/application/capture"
	"be-dashboard-nba/internal/presentation/response"

	"github.com/gofiber/fiber/v2"
)

func Proxy(service *captureApplication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status, contentType, body, err := service.Proxy(c.UserContext(), c.Body())
		if err != nil {
			if errors.Is(err, captureApplication.ErrInvalidPayload) {
				return response.Error(c, response.ErrorParam{Code: http.StatusBadRequest, Message: "Invalid capture JSON payload"})
			}
			if errors.Is(err, captureApplication.ErrBaseURLNotConfigured) {
				return response.Error(c, response.ErrorParam{Code: http.StatusInternalServerError, Message: "BASE_URL_CV is not configured"})
			}
			return response.Error(c, response.ErrorParam{Code: http.StatusBadGateway, Message: "Failed call CV service"})
		}

		c.Status(status)
		c.Set(fiber.HeaderContentType, contentType)
		return c.Send(body)
	}
}
