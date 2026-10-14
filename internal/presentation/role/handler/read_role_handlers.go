package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	req "be-dashboard-nba/internal/presentation/request"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/role/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get Roles List
// @Description  Retrieves a paginated and searchable list of all roles. Requires Bearer token.
// @Tags         Role
// @Produce      json
// @Param        search query    string false "Search term for role name"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  response.PaginatePayload{data=[]presenter.ReadRoleResponse} "Successfully get roles"
// @Failure      400  {object}  response.ErrorPayload "Invalid query parameters"
// @Failure      401  {object}  response.ErrorPayload "Unauthorized"
// @Failure      500  {object}  response.ErrorPayload "Failed get role"
// @Security     BearerAuth
// @Router       /roles [get]
func ReadRoles(svc contract.RoleUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.ReadRolesRequest
		if err = c.QueryParser(&request); err != nil {
			logger.Errorw("error parse request", err)
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadRolesUseCase(c.UserContext(), request)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed get role",
				})
		}

		return response.Paginate(c, req.PaginationPayload{
			Page:  data.Pagination.Page,
			Limit: data.Pagination.PageSize,
		}, int64(data.Pagination.TotalItems), response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadRoleListResponses(data.Data),
			Message: "Successfully get roles",
		})
	}
}


