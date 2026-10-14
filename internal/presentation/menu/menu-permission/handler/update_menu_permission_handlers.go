package handlers

import (
	"be-dashboard-nba/constant"
	contract "be-dashboard-nba/internal/domain/contract/usecase"
	"be-dashboard-nba/internal/infrastructure/logger"
	response "be-dashboard-nba/internal/presentation/response"
	"be-dashboard-nba/internal/infrastructure/validator"
	authInternal "be-dashboard-nba/internal/presentation/auth"
	"be-dashboard-nba/internal/presentation/menu/menu-permission/presenter"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// @Summary      Update Menu Permission
// @Description  Updates an existing menu permission (e.g., 'code', 'action_name') by its unique ID. Requires Bearer token.
// @Tags         Menu Permission
// @Accept       json
// @Produce      json
// @Param        menu_id            path      int                                     true "Menu ID"
// @Param        menu_permission_id path      int                                     true "Menu Permission ID"
// @Param        permission         body      presenter.UpdateMenuPermissionRequest true "Permission data to update"
// @Success      200 {object}  response.DataPayload           "Successfully update menu permission"
// @Failure      400 {object}  response.ErrorPayload        "Bad Request (Invalid ID, Invalid JSON, or Validation Error)"
// @Failure      401 {object}  response.ErrorPayload        "Unauthorized"
// @Failure      404 {object}  response.ErrorPayload        "Not Found (Permission ID not found)"
// @Failure      500 {object}  response.ErrorPayload        "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions/{menu_permission_id} [put]
func UpdateMenuPermission(svc contract.MenuPermissionUseCase, validate *validator.Validator) fiber.Handler {
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

		var payload presenter.UpdateMenuPermissionRequest
		if err = c.BodyParser(&payload); err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse request",
				})
		}

		if err := validate.Validate(payload); err != nil {
			return response.ErrorValidate(c, err)
		}

		ah, err := authInternal.GetAuth(c)
		if err != nil {
			logger.WithContext(c.UserContext()).Error(err, "error get auth handler")
			return
		}

		userID := ah.GetClaims().UserID

		err = svc.UpdateMenuPermissionUseCase(c.UserContext(), payload, userID, menuPermissionID)

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
			Message: "Successfully update menu permission",
		})
	}
}


