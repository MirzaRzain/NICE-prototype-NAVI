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

// @Summary      Delete Menu Permission
// @Description  Deletes a specific permission from a menu by its unique permission ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_id            path      int      true "Menu ID"
// @Param        menu_permission_id path      int      true "Menu Permission ID"
// @Success      200 {object}  response.DataPayload "Successfully delete menu permission"
// @Failure      400 {object}  response.ErrorPayload "Bad Request (Invalid ID)"
// @Failure      401 {object}  response.ErrorPayload "Unauthorized"
// @Failure      404 {object}  response.ErrorPayload "Not Found (Permission ID not found)"
// @Failure      500 {object}  response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions/{menu_permission_id} [delete]
func DeleteMenuPermission(svc contract.MenuPermissionUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuPermissionIDParams := c.Params("menu_permission_id")
		menuPermisionID, err := strconv.Atoi(menuPermissionIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		err = svc.DeleteMenuPermissionUseCase(c.UserContext(), menuPermisionID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuPermissionIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrMenuPermissionIdNotFound.Message,
					})
			} else {

				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusInternalServerError,
						Message: "Failed update menu permission",
					})
			}
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete menu permission",
		})
	}
}

