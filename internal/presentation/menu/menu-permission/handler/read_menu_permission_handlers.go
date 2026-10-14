package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	req "be-dashboard-nba/internal/presentation/request"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Get Menu Permission List
// @Description  Retrieves a paginated and searchable list of permissions (e.g., read, create) for a specific menu ID. Requires Bearer token.
// @Tags         Menu Permission
// @Produce      json
// @Param        menu_id  path      int    true   "ID of the Menu"
// @Param        search   query     string false  "Search term for action_name"
// @Param        page     query     int    false  "Page number (default: 1)"
// @Param        limit    query     int    false  "Items per page (default: 10)"
// @Param        order    query     string false  "Sort order (e.g., action_name ASC)"
// @Success      200      {object}  response.PaginatePayload{data=[]presenter.ReadMenuPermissionResponse} "Successfully get menu permission list"
// @Failure      400      {object}  response.ErrorPayload "Bad Request (Invalid Menu ID)"
// @Failure      401      {object}  response.ErrorPayload "Unauthorized"
// @Failure      404      {object}  response.ErrorPayload "Not Found (Menu ID not found)"
// @Failure      500      {object}  response.ErrorPayload "Internal Server Error (Query parsing or usecase failure)"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions [get]
func ReadMenuPermissionUseCase(svc contract.MenuPermissionUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "menu id params must be in number",
				})
		}

		var request presenter.ReadMenuPermissionListRequest
		if err = c.QueryParser(&request); err != nil {
			logger.Errorw("error parse request", err)
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadMenuPermissionUseCase(c.UserContext(), request, menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrMenuIdNotFound.Message,
					})
			} else {

				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusInternalServerError,
						Message: "Failed get menu permission",
					})
			}
		}
		return response.Paginate(c, req.PaginationPayload{
			Page:  data.Pagination.Page,
			Limit: data.Pagination.PageSize,
		}, int64(data.Pagination.TotalItems), response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuPermissionListResponses(data.Data),
			Message: "Successfully get menu permission list",
		})
	}
}


