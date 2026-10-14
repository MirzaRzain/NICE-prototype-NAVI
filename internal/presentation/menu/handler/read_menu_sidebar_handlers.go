package handlers

import (
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/menu/presenter"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// @Summary      Get User's Sidebar Menu
// @Description  Retrieves the menu list specifically for the authenticated user, filtered by their 'read' permissions and formatted in a nested tree structure for the sidebar.
// @Tags         Menu
// @Produce      json
// @Success      200  {object}  response.DataPayload{data=[]presenter.ReadMenuListResponse} "Successfully get sidebar menu"
// @Failure      401  {object}  response.ErrorPayload "Unauthorized"
// @Failure      500  {object}  response.ErrorPayload "Failed get sidebar menu"
// @Security     BearerAuth
// @Router       /menus/sidebar [get]
func ReadMenuSidebar(svc contract.MenuUseCase) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		data, err := svc.ReadSidebarMenuUseCase(c.UserContext(), userID)
		if err != nil {

			return response.Error(c, response.ErrorParam{
				Code:    http.StatusInternalServerError,
				Message: "Failed get sidebar menu",
			})
		}
		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuListResponse(data),
			Message: "Successfully get sidebar menu",
		})
	}
}


