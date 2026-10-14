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

// @Summary      Create Menu Permission
// @Description  Adds a new permission (e.g., 'read', 'create') to a specific menu by its ID. Requires Bearer token.
// @Tags         Menu Permission
// @Accept       json
// @Produce      json
// @Param        menu_id    path      int                                     true "Menu ID"
// @Param        permission body      presenter.CreateMenuPermissionRequest true "Permission data to create"
// @Success      201        {object}  response.DataPayload           "Successfully create menu permission"
// @Failure      400        {object}  response.ErrorPayload        "Bad Request (Invalid Menu ID, Invalid JSON, or Validation Error)"
// @Failure      401        {object}  response.ErrorPayload        "Unauthorized"
// @Failure      500        {object}  response.ErrorPayload        "Internal Server Error"
// @Security     BearerAuth
// @Router       /menus/{menu_id}/permissions [post]
func CreateMenuPermission(svc contract.MenuPermissionUseCase, validate *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		menuIDParams := c.Params("menu_id")
		menuID, err := strconv.Atoi(menuIDParams)
		if err != nil {
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusBadRequest,
					Message: "Failed parse params",
				})
		}

		var payload presenter.CreateMenuPermissionRequest
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

		err = svc.CreateMenuPermissionUseCase(c.UserContext(), payload, userID, menuID)

		if err != nil {
			if errors.Is(err, constant.ErrMenuIdNotFound) {
				return response.Error(c,
					response.ErrorParam{
						Code:    http.StatusNotFound,
						Message: constant.ErrDataNotFound.Message,
					})
			}
			return response.Error(c,
				response.ErrorParam{
					Code:    http.StatusInternalServerError,
					Message: "Failed create menu permission",
				})
		}

		return response.Data(c, response.DataPayload{
			Code:    http.StatusCreated,
			Data:    nil,
			Message: "Successfully create menu permission",
		})
	}
}


