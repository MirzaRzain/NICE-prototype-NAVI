package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/user/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get Current User Profile
// @Description  Fetches the profile details of the currently authenticated user. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Success      200 {object} response.DataPayload{data=presenter.ReadUserResponse} "Successfully retrieved user profile"
// @Failure      401 {object} response.ErrorPayload "Unauthorized (Invalid or missing token)"
// @Failure      404 {object} response.ErrorPayload "User profile not found"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/me [get]
func ReadProfileApp(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		ah, err := authInternal.GetAuth(c)
		if err != nil {
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusUnauthorized,
				Message: "Failed to get auth claims",
			})
		}

		data, err := svc.ReadUserProfileUseCase(c.UserContext(), ah.GetClaims().UserID)

		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    constant.ErrUserIdNotFound.Code,
					Message: constant.ErrUserIdNotFound.Message,
				})
			} else {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "failed to get user profile",
				})
			}
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadUserProfileUseCaseResponse(data),
			Message: "Successfully read user profile",
		})
	}
}

