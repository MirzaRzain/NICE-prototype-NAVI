package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Delete Role
// @Description  Deletes a role by its ID. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        role_id  path      int      true   "Role ID"
// @Success      200 {object}  response.DataPayload    "Successfully delete role"
// @Failure      400 {object}  response.ErrorPayload "Bad Request (Invalid Role ID)"
// @Failure      401 {object}  response.ErrorPayload "Unauthorized (Missing Auth)"
// @Failure      404 {object}  response.ErrorPayload "Role ID not found"
// @Failure      500 {object}  response.ErrorPayload "Failed delete role"
// @Security     BearerAuth
// @Router       /roles/{role_id} [delete]
func DeleteRole(svc contract.RoleUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		roleIDParams := c.Params("role_id")
		roleID, err := strconv.Atoi(roleIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		err = svc.DeleteRoleUseCase(c.UserContext(), roleID)

		if err != nil {
			if errors.Is(err, constant.ErrRoleIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrRoleIdNotFound.Message,
					})
			} else {

				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusInternalServerError,
						Message: "Failed update role",
					})
			}
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete role",
		})
	}
}

