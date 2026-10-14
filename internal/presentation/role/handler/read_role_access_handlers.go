package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	req "be-dashboard-nba/internal/presentation/request"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Get Role Menu Permission
// @Description  Retrieves all menu permissions associated with a specific role. Requires Bearer token.
// @Tags         Role Menu Permission
// @Produce      json
// @Param        role_id  path      int   true  "Role ID"
// @Success      200      {object}  response.PaginatePayload{data=presenter.RoleAccessResponse} "Successfully get role menu permission"
// @Failure      400      {object}  response.ErrorPayload "Invalid Role ID"
// @Failure      401      {object}  response.ErrorPayload "Unauthorized"
// @Failure      404      {object}  response.ErrorPayload "Role ID not found"
// @Failure      500      {object}  response.ErrorPayload "Failed get role menu permission"
// @Security     BearerAuth
// @Router       /roles/role-access/{role_id} [get]
func ReadRoleAccess(svc contract.RoleUseCase) fiber.Handler {
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

		var request presenter.ReadRoleAccessesRequest
		if err = c.QueryParser(&request); err != nil {
			logger.Errorw("error parse request", err)
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadRoleAccessUseCase(c.UserContext(), request, roleID)
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
						Message: "Failed get role menu permission",
					})
			}
		}
		return response.Paginate(c, req.PaginationPayload{
			Page:  data.Pagination.Page,
			Limit: data.Pagination.PageSize,
		}, int64(data.Pagination.TotalItems), response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadRoleAccessResponse(data.Data),
			Message: "Successfully get role detail",
		})
	}
}


