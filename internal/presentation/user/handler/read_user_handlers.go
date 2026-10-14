package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	req "be-dashboard-nba/internal/presentation/request"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/user/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// Summary      Get List of Users
// @Description  Retrieves a paginated and searchable list of users in the system. Requires Bearer token.
// @Tags         User
// @Produce      json
// @Param        search query    string false "Search term for name or email"
// @Param        page   query    int    false "Page number (default: 1)"
// @Param        limit  query    int    false "Items per page (default: 10)"
// @Param        order  query    string false "Sort order (e.g., name ASC)"
// @Success      200  {object}  response.PaginatePayload{data=[]presenter.ReadUserResponse} "Successfully get users"
// @Failure      400  {object}  response.ErrorPayload "Invalid query parameters"
// @Failure      401  {object}  response.ErrorPayload "Unauthorized"
// @Failure      500  {object}  response.ErrorPayload "Failed get users"
// @Security     BearerAuth
// @Router       /users [get]
func ReadUsers(svc contract.UserUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var request presenter.ReadUserRequest
		if err = c.QueryParser(&request); err != nil {
			logger.Errorw("error parse request", err)
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Invalid query parameters",
				})
		}

		data, err := svc.ReadUsersUseCase(c.UserContext(), request)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed get users",
				})
		}

		return response.Paginate(c, req.PaginationPayload{
			Page:  data.Pagination.Page,
			Limit: data.Pagination.PageSize,
		}, int64(data.Pagination.TotalItems), response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadUserResponses(data.Data),
			Message: "Successfully get users",
		})
	}
}


