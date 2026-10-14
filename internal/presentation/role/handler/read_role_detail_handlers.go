package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Get Role Detail
// @Description  Fetches the details of a single role by its ID. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        role_id  path      int      true   "Role ID"
// @Success      200 {object}  response.DataPayload{data=presenter.ReadRoleResponse} "Successfully get role detail"
// @Failure      400 {object}  response.ErrorPayload "Bad Request (Invalid Role ID)"
// @Failure      401 {object}  response.ErrorPayload "Unauthorized"
// @Failure      404 {object}  response.ErrorPayload "Not Found (Role ID not found)"
// @Failure      500 {object}  response.ErrorPayload "Failed get role detail"
// @Security     BearerAuth
// @Router       /roles/{role_id} [get]
func ReadRoleDetailUseCase(svc contract.RoleUseCase) fiber.Handler {
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

		data, err := svc.ReadRoleDetailUseCase(c.UserContext(), roleID)
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
						Message: "Failed get role detail",
					})
			}
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadRoleResponse(data),
			Message: "Successfully get role detail",
		})
	}
}

