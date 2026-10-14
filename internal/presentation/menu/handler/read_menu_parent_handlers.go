package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/menu/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get Parent Menus
// @Description  Retrieves a list of all root-level menus (menus without a parent), typically for use in a dropdown. Requires Bearer token.
// @Tags         Menu
// @Produce      json
// @Success      200 {object} response.DataPayload{data=[]presenter.MenuParent} "Successfully get parent menu"
// @Failure      401 {object} response.ErrorPayload "Unauthorized"
// @Failure      500 {object} response.ErrorPayload "Failed get parent menu"
// @Security     BearerAuth
// @Router       /menus/parent [get]
func ReadMenuParent(svc contract.MenuUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		data, err := svc.ReadMenuParentUseCase(c.UserContext())
		if err != nil {

			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed get parent menu",
			})
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuParentResponses(data),
			Message: "Successfully get parent menu",
		})
	}
}

