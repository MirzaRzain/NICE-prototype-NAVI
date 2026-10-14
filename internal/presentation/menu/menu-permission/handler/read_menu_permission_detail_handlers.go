package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Get Menu Permission Detail
// @Description  Fetches the details of a single menu permission by its unique ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_id            path      int      true "Menu ID"
// @Param        menu_permission_id path      int      true "Menu Permission ID"
// @Success      200 {object}  response.DataPayload{data=presenter.ReadMenuPermissionResponse} "Successfully get menu permission detail"
// @Failure      400 {object}  response.ErrorPayload "Bad Request (Invalid ID)"
// @Failure      401 {object}  response.ErrorPayload "Unauthorized"
// @Failure      404 {object}  response.ErrorPayload "Not Found (Permission ID not found)"
// @Failure      500 {object}  response.ErrorPayload "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions/{menu_permission_id} [get]
func ReadMenuPermissionDetailUseCase(svc contract.MenuPermissionUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuPermissionIDParams := c.Params("menu_permission_id")
		menuPermissionID, err := strconv.Atoi(menuPermissionIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})

		}

		data, err := svc.ReadMenuPermissionDetailUseCase(c.UserContext(), menuPermissionID)
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
						Message: "Failed get menu permission detail",
					})
			}
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuPermissionListResponse(data),
			Message: "Successfully get menu permission detail",
		})
	}
}

