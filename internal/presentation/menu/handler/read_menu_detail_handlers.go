package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/presentation/menu/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

func ReadMenuDetail(svc contract.MenuUseCase) fiber.Handler {
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

		data, err := svc.ReadMenuDetailUseCase(c.UserContext(), menuID)
		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusOK,
			Data:    presenter.ToReadMenuDetailResponse(data),
			Message: "Successfully get menu detail",
		})

	}
}

