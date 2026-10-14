package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/auth/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} response.DataPayload{data=presenter.UserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} response.ErrorPayload "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} response.ErrorPayload "User profile not found"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /auth/me [get]
func AuthMe(svc contract.AuthUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		data, err := svc.AuthMeUseCase(c.UserContext(), ah.GetClaims().UserID)

		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: constant.ErrUserIdNotFound.Error(),
				})
			}
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed logout",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadAuthMeResponse(data),
			Message: "Successfully read profile user",
		})
	}
}


