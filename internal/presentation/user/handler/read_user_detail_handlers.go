package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/user/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Summary      Get User Detail
// @Description  Fetches the details of a single user by their ID. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Param        user_id  path      string   true "User ID"
// @Success      200 {object} response.DataPayload{data=presenter.ReadUserDetailResponse} "Successfully get user detail"
// @Failure      400 {object} response.ErrorPayload "Bad Request (Invalid User ID)"
// @Failure      401 {object} response.ErrorPayload "Unauthorized"
// @Failure      404 {object} response.ErrorPayload "Not Found (User ID not found)"
// @Failure      500 {object} response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{user_id} [get]
func ReadUserDetail(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		userIDParams := c.Params("user_id")

		data, err := svc.ReadDetailUserUseCase(c.UserContext(), userIDParams)
		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    constant.ErrUserIdNotFound.Code,
					Message: constant.ErrUserIdNotFound.Message,
				})
			}

			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "failed to get profile",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadUserDetailResponse(data),
			Message: "Successfully read profile user",
		})
	}
}

