package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	response "be-dashboard-nba/internal/presentation/response"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      User Logout
// @Description  Logs out the current user by invalidating their session/token. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} response.DataPayload "Successfully logged out"
// @Failure      400 {object} response.ErrorPayload "Bad Request (e.g., Validation Error, Invalid JSON)"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error (e.g., Database connection issue)"
// @Security     BearerAuth
// @Router       /auth/logout [post]
func Logout(svc contract.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := authInternal.GetAuth(c)
		if err != nil {
			return
		}
		err = svc.LogoutUseCase(c.UserContext(), ah.GetClaims(), c.IP())
		if err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed logout",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully logout",
		})
	}
}

