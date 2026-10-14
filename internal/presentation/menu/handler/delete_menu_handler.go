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

// @Summary      Update Menu
// @Description  Updates an existing menu item by its ID. Requires Bearer token.
// @Tags         Menu
// @Accept       json
// @Produce      json
// @Param        menu_id  path      int                             true "Menu ID"
// @Success      200      {object}  response.DataPayload       "Successfully delete menu"
// @Failure      400      {object}  response.ErrorPayload  "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401      {object}  response.ErrorPayload  "Unauthorized"
// @Failure      500      {object}  response.ErrorPayload  "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id} [delete]
func DeleteMenu(svc contract.MenuUseCase) fiber.Handler {
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

		err = svc.DeleteMenuUseCase(c.UserContext(), menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrMenuIdNotFound.Message,
					})
			}
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    nil,
			Message: "Successfully delete menu",
		})

	}
}

