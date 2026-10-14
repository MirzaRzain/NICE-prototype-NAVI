package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Summary      Delete User
// @Description  Deletes a user by their ID. Requires Bearer token. Users cannot delete themselves.
// @Tags         User
// @Produce      json
// @Param        user_id path      string true "User ID"
// @Success      200     {object}  response.DataPayload "Successfully delete user"
// @Failure      400     {object}  response.ErrorPayload "Bad Request (Invalid User ID)"
// @Failure      401     {object}  response.ErrorPayload "Unauthorized"
// @Failure      403     {object}  response.ErrorPayload "Forbidden (Attempt to delete own account)"
// @Failure      404     {object}  response.ErrorPayload "Not Found (User ID not found)"
// @Failure      500     {object}  response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{user_id} [delete]
func DeleteUser(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		userIDParams := c.Params("user_id")

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.DeleteUserUseCase(c.UserContext(), userIDParams, userID)
		if err != nil {
			if errors.Is(err, constant.ErrUserIdNotFound) {
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusNotFound,
					Message: constant.ErrUserIdNotFound.Message,
				})
			}
			if errors.Is(err, constant.ErrForbiddenSelfDelete) {
				return response.Error(c, response.ErrorParam{
					Code:    constant.ErrForbiddenSelfDelete.Code,
					Message: constant.ErrForbiddenSelfDelete.Message,
				})
			}
			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed update user",
			})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete user",
		})
	}
}


