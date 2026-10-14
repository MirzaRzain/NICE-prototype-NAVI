package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/validator"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Update Role Menu Permission
// @Description  Updates menu permission assigned to a specific role. Requires Bearer token.
// @Tags         Role Menu Permission
// @Accept       json
// @Produce      json
// @Param        role_id  path      int                                     true  "Role ID"
// @Param        payload  body      presenter.UpdateRoleAccessRequest   true  "Role menu permission payload"
// @Success      200      {object}  response.DataPayload           "Successfully update role menu permission"
// @Failure      400      {object}  response.ErrorPayload          "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  response.ErrorPayload        "Unauthorized"
// @Failure      404      {object}  response.ErrorPayload        "Role ID or Menu Permission not found"
// @Failure      500      {object}  response.ErrorPayload        "Internal server error"
// @Security     BearerAuth
// @Router       /roles/role-access/{role_id} [put]
func UpdateRoleMenuPermission(svc contract.RoleUseCase, validate *validator.Validator) fiber.Handler {
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
		var request presenter.UpdateRoleAccessRequest
		fmt.Println(request)
		if err = c.BodyParser(&request); err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}
		if err := validate.Validate(request); err != nil {
			return response.ErrorValidate(c, err)
		}

		err = svc.UpdateRoleAccessUseCase(c.UserContext(), roleID, request)

		if err != nil {
			switch {
			case errors.Is(err, constant.ErrRoleIdNotFound):
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusNotFound,
					Message: constant.ErrRoleIdNotFound.Message,
				})
			case errors.Is(err, constant.ErrMenuPermissionIdNotFound):
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusNotFound,
					Message: constant.ErrMenuPermissionIdNotFound.Message,
				})
			case err.Error() == constant.ErrMsgValidate:
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: constant.ErrMsgValidate,
				})
			default:
				return response.Error(c, response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: constant.ErrMsgUnknownSource,
				})
			}
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully update role menu permission",
		})
	}
}

